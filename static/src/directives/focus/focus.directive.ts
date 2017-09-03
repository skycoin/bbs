import { Directive, Input, ElementRef, Renderer, OnChanges } from '@angular/core';

// tslint:disable-next-line:directive-selector
@Directive({ selector: '[focus]' })
export class FocusDirective implements OnChanges {
  @Input() focus: boolean;
  constructor(private el: ElementRef, private render: Renderer) { }

  ngOnChanges() {
    if (this.focus) {
      this.render.invokeElementMethod(this.el.nativeElement, 'focus');
    }
  }
}
