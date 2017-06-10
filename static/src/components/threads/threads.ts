import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService, Thread } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";

@Component({
  selector: 'threads',
  templateUrl: 'threads.html',
  styleUrls: ['threads.css'],
  encapsulation: ViewEncapsulation.None
})

export class ThreadsComponent implements OnInit {
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();
  constructor(private api: ApiService, private router: Router, private route: ActivatedRoute) { }
  threads: Array<Thread> = [];
  url: string = '';
  ngOnInit() {
    this.route.params.subscribe(res => {
      this.url = res['url'];
      this.start(res['board']);
    })
  }
  start(key) {
    this.api.getThreads(key).subscribe(threads => {
      this.threads = threads;
    });
  }

  open(master, ref: string) {
    this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
  }
}