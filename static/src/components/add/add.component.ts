import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'add',
  templateUrl: './add.component.html',
  styleUrls: ['./add.component.css']
})
export class AddComponent implements OnInit {
  select:string = 'board';
  form: {
    name?: string;
    description?: string;
    board?: string;
    thread?: string;
    seed?: string;
  } = { name: 'testName', description: 'tes', board: 'tes', thread: 'tes', seed: 'tess' }

  constructor() { }

  ngOnInit() {
  }
  add(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    console.log('form:', this.form);
  }
}
