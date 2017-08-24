export interface VotesSummary {
  up_votes?: number; // Total number of up-votes.
  down_votes?: number; // Total number of down-votes.
  current_user_voted?: boolean; // Whether current user has voted on this thread/post.
  current_user_vote_mode?: number; // (1: current user up-voted) | (-1: current user down-voted)
}

export interface Post {
  title?: string;
  body?: string;
  author?: string;
  created?: number;
  ref?: string;
  votes?: VotesSummary; // Posts now have vote summary here.
  uiOptions?: VoteOptions;
}

export interface VoteOptions {
  voted?: boolean;
  userVoted?: boolean;
  menu?: boolean;
}

export interface ThreadPage {
  board?: Board;
  thread?: Thread;
  posts?: Array<Post>;
}

export interface BoardPage extends Base {
  data?: BoardPageData;
}

export interface BoardPageData {
  board?: Board;
  threads?: Array<Thread>;
}

export interface Stats {
  node_is_master: boolean;
  node_cxo_address: string;
}

export interface Thread {
  title?: string;
  body?: string;
  created?: number;
  reference?: string;
  author_reference?: string;
  author_alias?: string;
  votes?: Votes;
  uiOptions?: VoteOptions;
}
export interface Votes {
  up?: VoteData;
  down?: VoteData;
  spam?: VoteData;
}
export interface VoteData {
  voted?: boolean;
  count?: number;
}
export interface Board {
  name?: string;
  description?: string;
  created?: number;
  submission_addresses?: Array<string>;
  public_key?: string;
  ui_options?: UIOptions; // custom param
}

export interface UIOptions {
  subscribe?: boolean;
}


export interface Subscription {
  synced?: boolean;
  accepted?: boolean;
  rejected_count?: number;
  config?: SubscriptionOption;
}

export interface SubscriptionOption {
  master?: boolean;
  public_key?: string;
  secret_key?: string;
}

export interface Base {
  okay?: boolean;
}

export interface AllBoardsData {
  master_boards?: Array<Board>;
  remote_boards?: Array<Board>;
}

export interface AllBoards extends Base {
  data?: AllBoardsData;
}
