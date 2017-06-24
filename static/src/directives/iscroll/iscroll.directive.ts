import { Directive, ElementRef, OnInit, OnDestroy } from '@angular/core';
import * as IScroll from 'iscroll';

@Directive({ selector: '[appIscroll]' })
export class IscrollDirective implements OnInit, OnDestroy {
  private iscroll: IScroll;
  constructor(private el: ElementRef) { }
  ngOnInit(): void {
    this.iscroll = new IScroll(this.el.nativeElement, { scrollbars: 'iScrollVerticalScrollbar' });
  }

  ngOnDestroy(): void {
    this.iscroll.destroy();
  }

}
