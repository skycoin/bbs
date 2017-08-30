export interface Connnections extends Base {
  data?: ConnnectionData;
}
export interface ConnnectionData {
  connections: Array<Connnection>;
}
export interface Connnection {
  address?: string;
  state?: string;
}
export interface Users extends Base {
  data: UserData;
}
export interface UserData {
  users?: Array<User>;
}
export interface User {
  alias?: string;
}
export interface VotesSummary {
  up_votes?: number; // Total number of up-votes.
  down_votes?: number; // Total number of down-votes.
  current_user_voted?: boolean; // Whether current user has voted on this thread/post.
  current_user_vote_mode?: number; // (1: current user up-voted) | (-1: current user down-voted)
}

export interface Post {
  name?: string;
  body?: string;
  creator?: string;
  created?: number;
  ref?: string;
  votes?: Votes;
  uiOptions?: VoteOptions;
  voteMenu?: boolean;
}

export interface VoteOptions {
  voted?: boolean;
  userVoted?: boolean;
  menu?: boolean;
}

export interface ThreadPage extends Base {
  data?: ThreadPageData;
}
export interface ThreadPageData {
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
  name?: string;
  body?: string;
  created?: number;
  ref?: string;
  creator?: string;
  // author_alias?: string;
  votes?: Votes;
  uiOptions?: VoteOptions;
}
export interface Votes {
  ref?: string;
  up_votes?: VoteData;
  down_votes?: VoteData;
}
export interface VoteData {
  voted?: boolean;
  count?: number;
}
export interface Board {
  name?: string;
  body?: string;
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
