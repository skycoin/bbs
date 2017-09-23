export interface AlertData {
  type?: string;
  title?: string;
  content?: string;
  footer?: string;
  duration?: number;
  autoDismiss?: boolean;
  clickEvent?: Function;
}
