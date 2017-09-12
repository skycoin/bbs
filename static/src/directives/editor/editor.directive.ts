import { Directive, ElementRef, forwardRef } from '@angular/core';
import * as Quill from 'quill';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

@Directive({
  // tslint:disable-next-line:directive-selector
  selector: '[editor]',
  providers: [{
    provide: NG_VALUE_ACCESSOR, useExisting:
    forwardRef(() => EditorDirective),
    multi: true
  }]
})
export class EditorDirective implements ControlValueAccessor {
  editor: any;
  toolbarOptions = [
    [{ 'header': [1, 2, 3, 4, 5, 6, false] }, { 'font': [] }],
    ['bold', 'italic', 'underline', 'strike', { 'color': [] }, { 'background': [] }, { 'script': 'sub' }, { 'script': 'super' }],
    ['blockquote', 'code-block', 'link', 'image'],
    [{ 'list': 'ordered' }, { 'list': 'bullet' }, { 'align': [] }]
  ];
  constructor(private el: ElementRef) {
    this.editor = new Quill(el.nativeElement, {
      theme: 'snow',
      modules: {
        toolbar: this.toolbarOptions
      }
    });
    this.editor.on('text-change', (delta, oldDelta, source) => {
      if (source === 'user' || source === 'api') {
        this.writeValue();
      }
    })
    const toolbar = this.editor.getModule('toolbar');
    toolbar.addHandler('image', (ev) => {
      if (ev) {
        const href = prompt('Enter the Image URL');
        if (href) {
          this.editor.insertEmbed(this.editor.getSelection(), 'image', href);
        }
      }
    });
  }
  onChange = (html: string) => { };

  onTouched = () => { };

  get value(): string {
    return this.editor.root.innerHTML;
  }
  writeValue(): void {
    this.onChange(this.value);
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
