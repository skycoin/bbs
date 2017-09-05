import { Component, OnInit, ViewEncapsulation, Input, HostBinding, HostListener } from '@angular/core';

@Component({
  selector: 'chip',
  template: `<i class="fa" [ngClass]="icon" aria-hidden="true" *ngIf="icon"></i>
  {{text}}
  <span *ngIf="show">{{count}}</span>`,
  styleUrls: ['./chip.component.scss'],
  encapsulation: ViewEncapsulation.None,
})

export class ChipComponent implements OnInit {
  @Input() color = '';
  @Input() text = '';
  @Input() count = 1;
  @Input() size = 'lg';
  @Input() show = true;
  @Input() icon = '';
  hoverTextColor = '';
  @HostBinding('style.borderColor') borderColor = '';
  @HostBinding('style.color') textColor = '';
  @HostBinding('style.backgroundColor') bgColor = '';
  constructor() { }

  ngOnInit() {
    this.init();
  }
  init() {
    if (this.count === undefined || this.count === null) {
      this.count = 0;
    }
    this.hoverTextColor = this.invertColor(this.color, true);
    if (this.color !== '') {
      this.borderColor = this.color;
      this.textColor = this.color;
    }
  }
  @HostListener('mouseenter')
  _mouseover() {
    if (this.show) {
      this.bgColor = this.color;
      this.textColor = this.hoverTextColor;
    }
  }
  @HostListener('mouseleave')
  _mouseleave() {
    if (this.show) {
      this.borderColor = this.color;
      this.textColor = this.color;
      this.bgColor = '';
    }
  }
  padZero(str: string, len?: number) {
    console.log('pad zero');
    len = len || 2;
    const zeros = new Array(len).join('0');
    return (zeros + str).slice(-len);
  }

  invertColor(hex: string, bw: boolean = true) {
    if (hex.indexOf('#') === 0) {
      hex = hex.slice(1);
    }
    // convert 3-digit hex to 6-digits.
    if (hex.length === 3) {
      hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
    }
    if (hex.length !== 6) {
      throw new Error('Invalid HEX color.');
    }
    const r = parseInt(hex.slice(0, 2), 16),
      g = parseInt(hex.slice(2, 4), 16),
      b = parseInt(hex.slice(4, 6), 16);
    if (bw) {
      return (r * 0.299 + g * 0.587 + b * 0.114) > 186
        ? '#000000'
        : '#FFFFFF';
    }
    return '#' + this.padZero((255 - r).toString(16)) + this.padZero((255 - g).toString(16)) + this.padZero((255 - b).toString(16));
  }
}
