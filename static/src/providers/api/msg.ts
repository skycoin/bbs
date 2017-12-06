export interface BoardBody {
  body?: string;
  created?: number;
  name?: string;
  submission_keys?: Array<string>;
  tags?: Array<string>;
  header?: BaseHeader;
}
export interface Board {
  body?: BoardBody;
  public_key?: string;
  ui_options?: UIOptions; // custom param
}

export interface AllBoardsData {
  master_boards?: Array<Board>;
  remote_boards?: Array<Board>;
}

export interface AllBoards extends Base {
  data?: AllBoardsData;
}

export interface LoginSessionUser extends User {
  public_key?: string;
  secret_key?: string;
}

export interface LoginSession {
  seed?: string;
  user?: LoginSessionUser;
}
export interface LoginData {
  logged_in?: boolean;
  session?: LoginSession;
}

export interface LoginInfo extends Base {
  data?: LoginData;
}

export interface FollowPage extends Base {
  data?: FollowPageData;
}
export interface FollowPageData {
  follow_page?: FollowPageDataInfo;
}
export interface FollowPageDataInfo {
  user_public_key?: string;
  following?: Array<FollowValue>;
  avoiding?: Array<FollowValue>;
}
export interface FollowValue {
  user_public_key?: string;
  tag?: string;
}
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
export interface PostBody {
  name?: string;
  creator?: string;
  body?: string;
  of_board?: string;
  of_thread?: string;
  ts?: number;
  type?: string;
}
export interface Post {
  body?: PostBody;
  header?: BaseHeader;
  votes?: Votes;
  uiOptions?: VoteOptions;
  voteMenu?: boolean;
  creatorMenu?: boolean;
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

export interface ThreadBody {
  body?: string;
  created?: number;
  creator?: string;
  name?: string;
  of_board?: string;
}
export interface Thread {
  body?: ThreadBody;
  name?: string;
  header?: BaseHeader;
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
export interface BaseHeader {
  type?: string;
  hash?: string;
  pk?: string;
  sig?: string;
}

export interface ThreadSubmission {
  name?: string;
  body?: string;
  created?: number;
  creator?; string;
  of_board?: string;
}

export interface PrepareRes extends Base {
  data: PrepareData
}

export interface PrepareData {
  hash?: string;
  raw?: string;
}
