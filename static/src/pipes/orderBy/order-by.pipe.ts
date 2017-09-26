import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'orderBy',
})
export class OrderByPipe implements PipeTransform {

  transform(values: Array<any>, args: string = 'desc'): Array<any> {
    if (!values || values.length < 1) {
      return [];
    }
    if (!values[0]) {
      return [];
    }
    if (!values[0].created) {
      return values;
    }
    switch (args) {
      case 'desc':
        values.sort((a, b) => {
          if (a.created > b.created) {
            return 1;
          }
          if (a.created < b.created) {
            return -1;
          }
          return 0;
        });
        break;
      case 'asc':
        values.sort((a, b) => {
          if (a.created > b.created) {
            return -1;
          }
          if (a.created < b.created) {
            return 1;
          }
          return 0;
        });
        break;
    }

    return values;
  }

}
