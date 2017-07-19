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
  author?: string;
  created?: string;
  master_board?: string;
  ref?: string;
  votes?: VotesSummary; // Threads now have vote summary here.
  uiOptions?: VoteOptions;
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
