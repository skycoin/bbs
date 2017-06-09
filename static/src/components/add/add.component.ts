import { Component, OnInit } from '@angular/core';
import { ApiService } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";

@Component({
  selector: 'add',
  templateUrl: './add.component.html',
  styleUrls: ['./add.component.css']
})
export class AddComponent implements OnInit {
  select: string = 'board';
  form: {
    name?: string;
    description?: string;
    board?: string;
    thread?: string;
    seed?: string;
    title?: string;
    body?: string;
    fromBoard?: string;
    toBoard?: string;
  } = {
    name: '',
    description: '',
    board: '',
    thread: '',
    seed: '',
    title: '',
    body: '',
    fromBoard: '',
    toBoard: ''
  }

  constructor(private api: ApiService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit() {
    this.route.params.subscribe(data => {
      if (data['exec']) {
        this.select = data['exec'];
      }
      this.form.board = data['board'];
      this.form.thread = data['thread'];
    });
  }
  init() {
    this.form = {
      name: '',
      description: '',
      board: '',
      thread: '',
      seed: '',
      title: '',
      body: '',
      fromBoard: '',
      toBoard: ''
    }
  }
  clear(ev) {
    this.select = ev;
    this.init();
  }
  add(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    let data = new FormData();
    // console.log('form:', this.form);
    switch (this.select) {
      case 'board':
        data.append('name', this.form.name);
        data.append('description', this.form.description);
        data.append('seed', this.form.seed);
        this.api.addBoard(data).then(res => {
          alert('add success');
          this.init();
        }, err => {
          console.error('err:', err);
        });
        break;
      case 'thread':
        data.append('board', this.form.board);
        data.append('description', this.form.description);
        data.append('name', this.form.name);
        this.api.addThread(data).then(res => {
          alert('add success');
          this.init();
          console.log('add success:', res);
        }, err => {
          console.error('err:', err);
        });
        break;
      case 'post':
        data.append('board', this.form.board);
        data.append('thread', this.form.thread);
        data.append('title', this.form.title);
        data.append('body', this.form.body);
        this.api.addPost(data).then(res => {
          console.log('add success:', res);
          this.init();
          alert('add success');
        }, err => {
          console.error('err:', err);
        });
        break;
      case 'changeBoard':
        data.append('from_board', this.form.fromBoard);
        data.append('to_board', this.form.toBoard);
        data.append('thread', this.form.thread);
        this.api.changeThread(data).then(res => {
          console.log('add success:', res);
          this.init();
          alert('add success');
        }, err => {
          console.error('err:', err);
        });
        break;
    }
  }
}
