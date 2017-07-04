import {Directive, ElementRef, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';

import * as Clipboard from 'clipboard';

@Directive({
    selector: '[appClip]'
})
export class ClipDirective implements OnInit, OnDestroy {
    private clipBoard: Clipboard;
    @Input() clipText = '';
    @Output() onClip: EventEmitter<boolean> = new EventEmitter<boolean>();

    constructor(private el: ElementRef) {
    }

    ngOnInit() {
        const option: Clipboard.Options = {
            text: (ele) => {
                return this.clipText;
            }
        };
        this.clipBoard = new Clipboard(this.el.nativeElement, option);
        this.clipBoard.on('success', () => {
            this.onClip.emit(true);
        });
        this.clipBoard.on('error', () => {
            this.onClip.emit(false);
        });
    }

    ngOnDestroy() {
        if (this.clipBoard) {
            this.clipBoard.destroy();
        }
    }
}
