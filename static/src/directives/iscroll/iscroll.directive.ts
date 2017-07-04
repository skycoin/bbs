import {Directive, ElementRef, OnDestroy, OnInit} from '@angular/core';
import * as IScroll from 'iscroll';

@Directive({selector: '[appIscroll]'})
export class IscrollDirective implements OnInit, OnDestroy {
    private iscroll;

    constructor(private el: ElementRef) {
    }

    ngOnInit(): void {
        this.iscroll = new IScroll(this.el.nativeElement, {scrollbars: 'iScrollVerticalScrollbar'});
    }

    ngOnDestroy(): void {
        this.iscroll.destroy();
    }

}
