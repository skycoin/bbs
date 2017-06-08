import { Component, OnInit } from '@angular/core';
import { ApiService } from "../../providers";
@Component({
  selector: 'threadPage',
  templateUrl: 'threadPage.html',
  styleUrls: ['threadPage.css']
})

export class ThreadPageComponent implements OnInit {
  data: { posts: Array<any>, thread: any } = { posts: [], thread: { name: '', description: '' } };
  constructor(private api: ApiService) { }

  ngOnInit() { }

  open(master, ref: string) {
    console.warn('open:', master);
    this.api.getThreadpage(master, ref).then(data => {
      console.warn('get threads2:', data);
      this.data = <{ posts: Array<any>, thread: any }>data;
      console.log('this data:', this.data);
    });
  }
}