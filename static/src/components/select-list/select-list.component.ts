import { Component, OnInit, forwardRef, Input } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

@Component({
  selector: 'select-list',
  templateUrl: './select-list.component.html',
  styleUrls: ['./select-list.component.scss'],
  providers: [{
    provide: NG_VALUE_ACCESSOR, useExisting:
    forwardRef(() => SelectListComponent),
    multi: true
  }]
})

export class SelectListComponent implements OnInit, ControlValueAccessor {
  showList = false;
  @Input() list = [];
  selectValue = '';
  constructor() { }

  ngOnInit() { }
  select(ev: Event, value: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.selectValue = value;
    this.showList = false;
    this.writeValue(value);
  }
  onChange = (value: string) => { };

  onTouched = () => { };

  get value(): string {
    return this.selectValue;
  }
  writeValue(value: string): void {
    this.onChange(value);
  }
  registerOnChange(fn: (html: string) => void): void {
    this.onChange = fn;
  }
  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }
  setDisabledState(isDisabled: boolean): void {

  }
}
