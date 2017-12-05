import { Directive, ElementRef, Input, OnInit, HostListener } from '@angular/core';
import { UserService, ConnectServiceInfo } from '../../service';


@Directive({
  selector: '[appLabel]'
})
export class LabelDirective implements OnInit {
  @Input() text = '';
  @Input() info: ConnectServiceInfo | null;
  @Input() index = -1;
  @Input() key = '';
  isEdit = false;
  constructor(private el: ElementRef, private user: UserService) { }
  ngOnInit() {
    this.el.nativeElement.value = this.text;
    this.change();
  }
  change() {
    if (this.el.nativeElement.value.length > 6) {
      this.el.nativeElement.value = this.el.nativeElement.value.substring(0, 5) + '...';
    }
  }
  @HostListener('focus', ['$event'])
  _foucs(ev: Event) {
    this.isEdit = true;
    this.el.nativeElement.value = this.text;
  }
  @HostListener('blur', ['$event'])
  _blur(ev: Event) {
    this.isEdit = false;
    const value = this.el.nativeElement.value;
    if (this.text !== value) {
      this.info.label = value;
      this.user.editClientConnectInfo(this.info, this.key, this.index);
    }
    this.change();
  }
}
