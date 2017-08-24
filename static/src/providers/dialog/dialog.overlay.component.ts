import { Component, OnInit, ViewEncapsulation, ComponentRef, HostListener } from '@angular/core';

@Component({
  selector: 'app-dialog-overlay',
  template: '',
  styles: [`
  app-dialog-overlay {
    position:fixed;
    top:0;bottom:0;left:0;right:0;
    background-color: rgba(0, 0, 0, 0.4);
    z-index: 2020;
  }
  `],
  encapsulation: ViewEncapsulation.None
})

export class DialogOverlayComponent implements OnInit {
  self: ComponentRef<DialogOverlayComponent>;
  constructor() { }

  ngOnInit() { }

  @HostListener('click', ['$event'])
  _click(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const overlayNativeEl = this.self.location.nativeElement;
    overlayNativeEl.parentNode.removeChild(overlayNativeEl);
    this.self.destroy();
    this.self = null;
  }
}
