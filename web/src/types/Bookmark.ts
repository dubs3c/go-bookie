

export interface Bookmark {
    id: number;
    title: string;
    url: string;
    body: string;
    description: string;
    image: string;
    archived: Boolean;
    deleted: Boolean;
    tags: string
};


export interface BookmarkList {
    id: string;
    title: string;
    url: string;
    description: string;
    image: string;
    archived: Boolean;
    deleted: Boolean;
    tags: string
};