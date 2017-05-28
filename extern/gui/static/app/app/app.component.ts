import {Component} from '@angular/core';

@Component({
    selector: 'pm-app',
    template:`
    <boards-list></boards-list>
    `
})

export class AppComponent{
    pageTitle: string = `Acme Product Management`; 
}
