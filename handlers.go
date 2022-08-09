package gobookie

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/netip"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	goose "github.com/advancedlogic/GoOse"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
)

// Server struct
type Server struct {
	DB     *pgxpool.Pool
	Router *chi.Mux
}

func getarticle(ctx context.Context, db *pgxpool.Pool, id int, url string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()
	g := goose.New()
	article, err := g.ExtractFromURL(url)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, "UPDATE bookmarks SET body=$2, title=$3, description=$4, image=$5 WHERE id=$1", id, article.RawHTML, article.Title, article.MetaDescription, article.TopImage)
	if err != nil {
		cancel()
		return err
	}
	return nil
}

// isValidURL checks if submitted link is an URL. Does not allow IPs or basic auth.
func isValidURL(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}

	if u.IsAbs() && u.Scheme == "https" || u.Scheme == "http" && u.Hostname() != "localhost" && u.User.String() == "" {
		_, err := netip.ParseAddr(u.Hostname())
		if err == nil {
			return false
		}

		return true
	}

	return false
}

// Only admin can run user APIs
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &CreateUserRequest{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		RespondWithError(w, 400, "Could not decode json")
	}

	if user.Password != user.VerifyPassword {
		RespondWithError(w, 400, "Passwords does not match")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		RespondWithError(w, 400, "Email is not valid")
	}

	_, err = s.UserRepositoryCreate(user)

	if err != nil {
		log.Println(err)
		RespondWithError(w, 500, "Could not create user")
	}

	RespondWithStatusCode(w, 201)

}

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := s.UserRepositoryList()

	if err != nil {
		log.Println(err)
		RespondWithError(w, 500, "Could not list users")
		return
	}

	RespondWithJSON(w, 200, &users)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {

	// should check here if user has a cookie with correct token
	// to avoid further processing

	u := &UserLogin{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		RespondWithError(w, 400, "Could not decode json")
		return
	}

	userId, err := s.UserRepositoryExists(u.Email, u.Password)

	if err != nil {
		log.Println(err)
		if userId == 0 {
			RespondWithError(w, 404, "Username or Password Incorrect")
			return
		}
		RespondWithError(w, 500, "Could not attempt login")
		return
	}

	if userId > 0 {
		token, err := s.UserRepositoryLogin(userId)
		if err != nil {
			RespondWithError(w, 500, "Could not login")
			return
		}

		// set secure flag in prod

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			MaxAge:   604800,
			HttpOnly: true,
			Path:     "/",
		})

		//w.Header().Add("Set-Cookie", "token="+token+"; path=/; HttpOnly; SameSite=Lax; Max-Age=604800;")

		RespondWithStatusCode(w, 200)
	} else {
		RespondWithStatusCode(w, 404)
	}

}

func (s *Server) UserLogout(w http.ResponseWriter, r *http.Request) {

}

// CreateBookmark Create bookmark
func (s *Server) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	data := &CreateBookmarkRequest{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		RespondWithError(w, 400, "Could not decode json")
	}

	if data.URL == "" {
		RespondWithError(w, 400, "You need to enter a URL")
		return
	}

	LastInsertedID, err := s.BookmarkRepositoryInsert(data)

	go func() {
		if isValidURL(data.URL) {
			err := getarticle(context.Background(), s.DB, LastInsertedID, data.URL)
			if err != nil {
				log.Printf("ERROR: Could not fetch website %s: %s", data.URL, err.Error())
			}
		}
	}()

	if err != nil {
		RespondWithError(w, 500, "Could not create bookmark")
		return
	}

	if LastInsertedID == 0 {
		RespondWithError(w, 500, "Bookmark was not created")
		return
	}

	RespondWithStatusCode(w, 201)
}

// ListBookmarks Return all bookmarks
func (s *Server) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	var (
		pageNumber    int
		bookmarkLimit int
		total         int
		archived      bool
		deleted       bool
		tags          string
		err           error
		tagList       []string
	)

	if pageNumber, err = strconv.Atoi(r.URL.Query().Get("page")); err != nil {
		pageNumber = 1
	}

	if bookmarkLimit, err = strconv.Atoi(r.URL.Query().Get("limit")); err != nil {
		bookmarkLimit = 10
	}

	if archived, err = strconv.ParseBool(r.URL.Query().Get("archived")); err != nil {
		archived = false
	}

	if deleted, err = strconv.ParseBool(r.URL.Query().Get("deleted")); err != nil {
		deleted = false
	}

	if tags = r.URL.Query().Get("tags"); tags != "" {
		tagList = strings.Split(tags, ",")
	} else {
		tagList = []string{}
	}

	bookmarks, err := s.BookmarkRepositoryGetAllBookmarks(pageNumber, bookmarkLimit, archived, deleted, tagList)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong while fetching bookmarks")
		log.Println(err)
		return
	}

	if len(bookmarks) < bookmarkLimit {
		total = 1
	} else {
		bookmarksAmount, _ := s.BookmarkRepositoryCount()
		total = (bookmarksAmount + bookmarkLimit - 1) / bookmarkLimit
	}

	paginatedResult := PaginatedBookmarks{
		Page:       pageNumber,
		TotalPages: total,
		Limit:      bookmarkLimit,
		Data:       bookmarks,
	}

	RespondWithJSON(w, 200, paginatedResult)
}

// GetBookmark Get a specific bookmark by its ID
func (s *Server) GetBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")
	var htmlBody bool
	var err error

	if htmlBody, err = strconv.ParseBool(r.URL.Query().Get("htmlbody")); err != nil {
		htmlBody = false
	}

	bookmark, err := s.BookmarkRepositoryGetBookmarkByID(bookmarkID)

	if err != nil {
		log.Println(err)
		RespondWithError(w, 500, "Could not fetch bookmark")
		return
	}

	if bookmark.ID < 1 {
		RespondWithStatusCode(w, 404)
		return
	}

	if htmlBody {
		w.Write([]byte(bookmark.Body))
	} else {
		RespondWithJSON(w, 200, bookmark)
	}
}

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryDeleteBookmarkByID(bookmarkID)

	if err != nil {
		RespondWithError(w, 500, "Could not delete bookmark")
		return
	}
	RespondWithJSON(w, 200, bookmark)
}

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) MoveBookmarkToTrash(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryTrashBookmarkByID(bookmarkID)

	if err != nil {
		RespondWithError(w, 500, "Could not mark bookmark as deleted")
		return
	}
	RespondWithJSON(w, 200, bookmark)
}

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) CorsOptions(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, 200, "")
}

// UpdateBookmark Update a bookmark by its ID
func (s *Server) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	var data UpdateBookmarkRequest
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		RespondWithError(w, 400, "Could not read json")
		return
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		RespondWithError(w, 400, "Could not decode json")
		return
	}

	bookmark, err := s.BookmarkRepositoryGetBookmarkByID(bookmarkID)

	if err != nil {
		RespondWithError(w, 500, "Could not fetch bookmark")
		return
	}

	// Because PUT and PATCH performs identical operations,
	// this is an easy way to perform different database
	// operations based on the method.
	switch r.Method {
	case "PUT":

		if data.Archived != nil {
			bookmark.Archived = *data.Archived
		}
		if data.Deleted != nil {
			bookmark.Deleted = *data.Deleted
		}
		bookmark.Body = data.Body
		bookmark.Description = data.Description
		bookmark.Image = data.Image
		bookmark.Title = data.Title
		bookmark.URL = data.URL

		if err = s.BookmarkRepositoryUpdateBookmark(bookmark); err != nil {
			RespondWithError(w, 500, "Could not update bookmark")
			return
		}
	case "PATCH":
		err := patch(data, &bookmark)
		if err != nil {
			RespondWithError(w, 500, "Something went wrong when trying to patch bookmark...")
			return
		}
		if err = s.BookmarkRepositoryPatchBookmark(bookmark); err != nil {
			RespondWithError(w, 500, "Could not patch bookmark")
			return
		}
	default:
		RespondWithStatusCode(w, 400)
		return
	}

	RespondWithStatusCode(w, 200)
}

func (s *Server) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := s.TagsRepositoryGetTags()

	if err != nil {
		RespondWithError(w, 500, "Could not fetch tags")
		return
	}
	RespondWithJSON(w, 200, tags.Tags)
}

func (s *Server) CreateTag(w http.ResponseWriter, r *http.Request) {
	tagData := &CreateTagRequest{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(tagData); err != nil {
		RespondWithError(w, 400, "Error decoding json. The bookmark id needs to be an interger and the tag name a string")
	}

	if tagData.TagName == "" {
		RespondWithError(w, 400, "You need to specify a tag name")
		return
	}

	if err := s.TagsRepositoryCreateTag(tagData.BookmarkID, tagData.TagName); err != nil {
		RespondWithError(w, 500, "Could not create tag")
		log.Println(err)
		return
	}

	RespondWithStatusCode(w, 201)
}

func (s *Server) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagData := &CreateTagRequest{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(tagData); err != nil {
		RespondWithError(w, 400, "Error decoding json. The bookmark id needs to be an interger and the tag name a string")
	}

	if tagData.TagName == "" {
		RespondWithError(w, 400, "You need to specify a tag name")
		return
	}

	if err := s.TagsRepositoryDeleteTagByBookmarkIDAndTagID(tagData.BookmarkID, tagData.TagName); err != nil {
		log.Println(err)
		RespondWithError(w, 500, "Could not delete tag")
		return
	}

	RespondWithStatusCode(w, 204)
}

func (s *Server) UpdateTag(w http.ResponseWriter, r *http.Request) {
	// because we have a many-2-many relationship, a user should not be able to update a tag name
	// if other users have the same tag. Instead, a new tag should be created with the updated name.
	// Consider investigating if creating a SQL trigger/function could contain this logic.
	// IF count(tag.id) > 1 THEN create_new_tag() ELSE update_tag()
}

// Error - A function that returns a custom error message
func Error(msg string) error {
	return errors.New(msg)
}

// Patch - Update a destination struct with non empty values from a src struct
// This function looks at src and takes the non empty values
// and assigns it to dst on the correct field
func patch(v interface{}, bookmark *Bookmark) error {
	src := reflect.ValueOf(v)
	if !src.IsValid() {
		return Error("src reflect not valid")
	}

	if src.Type().Kind() != reflect.Struct {
		return Error("src parameter is not a struct")
	}

	dst := reflect.ValueOf(bookmark).Elem()
	if !dst.IsValid() {
		return Error("dst reflect not valid")
	}

	for i := 0; i < src.NumField(); i++ {
		if src.Field(i).Kind() == reflect.String {
			if src.Field(i).String() != "" {
				fieldName := src.Type().Field(i).Name
				x := dst.FieldByName(fieldName)
				if x.IsValid() {
					if !x.CanSet() {
						return Error("ERROR: Can not update destintion struct")
					}
					x.SetString(src.Field(i).String())
				}
			}
		}

		if src.Field(i).Kind() == reflect.Ptr {
			l := src.Field(i).Elem()
			if l.IsValid() {
				fieldPointer := reflect.Indirect(src.Field(i))
				if fieldPointer.Type().Kind() == reflect.Bool {
					fieldName := src.Type().Field(i).Name
					x := dst.FieldByName(fieldName)

					if !x.CanSet() {
						return Error("ERROR: Can not update destintion struct")
					}

					x.SetBool(fieldPointer.Bool())
				}
			}
		}
	}
	return nil
}
