import { Component, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-loading',
  templateUrl: './loading.html',
  styleUrls: ['./loading.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class LoadingComponent {
  loadingText = 'Loading';
  iconColor = ''
  textColor = ''
  constructor() {
  }
}
