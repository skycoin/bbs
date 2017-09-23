import { Component, OnInit, ViewEncapsulation, HostListener, HostBinding } from '@angular/core';

@Component({
  selector: 'to-top',
  template: `<a href="javascript:void(0);" class="to-top" (click)="toTop($event)">
  <i class="fa fa-chevron-up" aria-hidden="true"></i></a>`,
  styleUrls: ['./to-top.component.scss'],
  encapsulation: ViewEncapsulation.None,
})

export class ToTopComponent implements OnInit {
  @HostBinding('style.display') show = 'none';
  constructor() { }

  ngOnInit() { }
  toTop(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.scrollToTop(500);
  }

  scrollToTop(scrollDuration) {
    const scrollHeight = window.scrollY,
      scrollStep = Math.PI / (scrollDuration / 15),
      cosParameter = scrollHeight / 2;
    let scrollCount = 0;
    const scrollInterval = setInterval(function () {
      if (window.scrollY !== 0) {
        scrollCount = scrollCount + 1;
        const scrollMargin = cosParameter - cosParameter * Math.cos(scrollCount * scrollStep);
        window.scrollTo(0, (scrollHeight - scrollMargin));
      } else {
        clearInterval(scrollInterval);
      }
    }, 15);
  }

  @HostListener('window:scroll')
  windowScroll() {
    const pos = (document.documentElement.scrollTop || document.body.scrollTop) + document.documentElement.offsetHeight;
    const max = document.documentElement.scrollHeight;
    const clientHeight = document.documentElement.clientHeight;
    const distance = max - pos;
    const enableScroll = max - clientHeight - 10;
    if (distance < (enableScroll - 20)) {
      this.show = 'block';
    } else {
      this.show = 'none';
    }
  }
}
