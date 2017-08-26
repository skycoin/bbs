import { Component, OnInit, ViewEncapsulation, Input, Output, EventEmitter, OnChanges, SimpleChanges } from '@angular/core';

@Component({
  selector: 'vote-box',
  templateUrl: 'vote-box.component.html',
  styleUrls: ['./vote-box.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class VoteBoxComponent implements OnInit, OnChanges {
  @Input() up = 0
  @Input() down = 0;
  upText = '';
  downText = '';
  @Output() onUp = new EventEmitter<any>();
  @Output() onDown = new EventEmitter<any>();
  constructor() { }

  ngOnInit() {
  }
  ngOnChanges(changes: SimpleChanges): void {
    if (this.up !== null && this.down !== null) {
      const total = this.up + this.down;
      this.upText = ((this.up / total) * 100) + '%';
      this.downText = ((this.down / total) * 100) + '%';
    }
  }
  actionUp(ev: Event) {
    this.onUp.emit();
  }
  actionDown(ev: Event) {
    this.onDown.emit();
  }
}
