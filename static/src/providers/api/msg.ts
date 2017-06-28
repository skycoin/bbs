export interface Post {
  title?: string;
  body?: string;
  author?: string;
  created?: number;
  ref?: string;
}

export interface ThreadPage {
  board?: Board;
  thread?: Thread;
  posts?: Array<Post>;
}

export interface BoardPage {
  board?: Board;
  threads?: Array<Thread>;
}

export interface Stats {
  node_is_master: boolean;
  node_cxo_address: string;
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
  address?: Array<string>;
  created?: number;
  ui_options?: UIOptions;
}

export interface UIOptions {
  subscribe?: boolean;
}



export interface SubScription {
  synced?: boolean;
  accepted?: boolean;
  rejected_count?: number;
  config?: SubScriptionOptions;
}

export interface SubScriptionOptions {
  master?: boolean;
  public_key?: string;
  secret_key?: string;
}
