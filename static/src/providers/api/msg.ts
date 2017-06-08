export interface Post {
  title?: string;
  body?: string;
  author?: string;
  created?: number;
  ref?: string;
}

export interface ThreadPage {
  thread?: Thread;
  posts?: Array<Post>;
}

export interface Thread {
  name?: string;
  description?: string;
  master_board?: string;
  ref?: string;
}

export interface Board {
  name?: string;
  description?: string;
  public_key?: string;
  url?: string;
  created?: number;
}