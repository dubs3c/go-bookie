import type {BookmarkList} from "./Bookmark"

export interface Pagination {
    page: number;
    totalPages: number;
    limit: number;
    data: BookmarkList[];
  }

