import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'cutText'
})

export class CutTextPipe implements PipeTransform {
  transform(value: any, length: number = 15): string {
    if (!value) {
      return '';
    }
    if (value.length < length) {
      return value;
    } else {
      return value.substring(0, length) + '...';
    }
  }
}
