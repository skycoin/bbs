import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService } from "../../providers";
@Component({
  selector: 'threads',
  templateUrl: 'threads.html',
  styleUrls: ['threads.css'],
  encapsulation: ViewEncapsulation.None
})

export class ThreadsComponent implements OnInit {
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();
  constructor(private api: ApiService) { }
  threads: Array<any> = [];
  ngOnInit() {
  }
  start(key) {
    this.api.getThreads(key).then(data => {
      console.warn('get threads:', data);
      this.threads = <Array<any>>data;
    });
  }

  open(master, ref: string) {
    this.thread.emit({ master: master, ref: ref });
  }
}