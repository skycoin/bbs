import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'replace'
})
/**
 * example
 * {{data | replace: [regexpStr,replaceStr]}}
 */
export class RepalcePipe implements PipeTransform {
  transform(value: string, ...args: string[]): string {
    const regexpStr = args[0][0];
    const replaceStr = args[0][1];
    return value.replace(regexpStr, replaceStr);
  }
}
