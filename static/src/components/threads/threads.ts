import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService } from "../../providers";
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
  threads: Array<any> = [];
  ngOnInit() {
    this.route.params.subscribe(res => {
      this.start(res['board']);
    })
  }
  start(key) {
    this.api.getThreads(key).then(data => {
      console.warn('get threads:', data);
      this.threads = <Array<any>>data;
    });
  }

  open(master, ref: string) {
    this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
    // this.thread.emit({ master: master, ref: ref });
  }
}