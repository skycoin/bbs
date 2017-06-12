import { Component, OnInit } from '@angular/core';
import { ApiService, ThreadPage } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";
import { FormControl, FormGroup } from '@angular/forms';
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'threadPage',
  templateUrl: 'threadPage.html',
  styleUrls: ['threadPage.css']
})

export class ThreadPageComponent implements OnInit {
  boardKey: string = '';
  threadKey: string = '';
  data: ThreadPage = { posts: [], thread: { name: '', description: '' } };
  postForm = new FormGroup({
    title: new FormControl(),
    body: new FormControl()
  });
  constructor(private api: ApiService, private router: Router, private route: ActivatedRoute, private modal: NgbModal) { }

  ngOnInit() {
    this.route.params.subscribe(res => {
      this.boardKey = res['board'];
      this.threadKey = res['thread'];
      this.open(this.boardKey, this.threadKey);
    })
  }

  openReply(content) {
    this.modal.open(content).result.then((result) => {
      let data = new FormData();
      data.append('board', this.boardKey);
      data.append('thread', this.threadKey);
      data.append('title', this.postForm.get('title').value);
      data.append('body', this.postForm.get('body').value);
      this.api.addPost(data).subscribe(post => {
        if (post) {
          this.data.posts.unshift(post);
        }
      })
    },err => {});

  }

  reply() {
    if (!this.boardKey || !this.threadKey) {
      alert('Will not be able to post');
      return;
    }
    this.router.navigate(['/add', { exec: 'post', board: this.boardKey, thread: this.threadKey }]);
  }
  open(master, ref: string) {
    console.warn('open:', master);
    this.api.getThreadpage(master, ref).subscribe(data => {
      this.data = data;
    });
  }
}