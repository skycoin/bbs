import { Component, OnInit } from '@angular/core';
import { ApiService } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";

@Component({
  selector: 'threadPage',
  templateUrl: 'threadPage.html',
  styleUrls: ['threadPage.css']
})

export class ThreadPageComponent implements OnInit {
  boardKey: string = '';
  threadKey: string = '';
  data: { posts: Array<any>, thread: any } = { posts: [], thread: { name: '', description: '' } };
  constructor(private api: ApiService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit() {
    this.route.params.subscribe(res => {
      this.boardKey = res['board'];
      this.threadKey = res['thread'];
      this.open(this.boardKey, this.threadKey);
    })
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
    this.api.getThreadpage(master, ref).then(data => {
      console.warn('get threads2:', data);
      this.data = <{ posts: Array<any>, thread: any }>data;
      console.log('this data:', this.data);
    });
  }
}