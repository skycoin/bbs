import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService, Thread } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";
import { FormControl, FormGroup } from '@angular/forms';
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'threads',
  templateUrl: 'threads.html',
  styleUrls: ['threads.css'],
  encapsulation: ViewEncapsulation.None
})

export class ThreadsComponent implements OnInit {
  @Output() thread: EventEmitter<{ master: string, ref: string }> = new EventEmitter();
  constructor(private api: ApiService, private router: Router, private route: ActivatedRoute, private modal: NgbModal) { }
  threads: Array<Thread> = [];
  boardKey:string = '';
  url: string = '';
  addForm = new FormGroup({
    description: new FormControl(),
    name: new FormControl()
  });
  ngOnInit() {
    this.route.params.subscribe(res => {
      this.url = res['url'];
      this.boardKey = res['board'];
      this.start(res['board']);
    })
  }
  start(key) {
    this.api.getThreads(key).subscribe(threads => {
      this.threads = threads;
    });
  }
  openAdd(content) {
    this.modal.open(content).result.then((result) => {
      if (result) {
        let data = new FormData();
        data.append('board',this.boardKey);
        data.append('description', this.addForm.get('description').value);
        data.append('name', this.addForm.get('name').value);
        this.api.addThread(data).subscribe(thread => {
          this.threads.unshift(thread);
        });
      }
    },err => {});
  }
  open(master, ref: string) {
    this.router.navigate(['p', { board: master, thread: ref }], { relativeTo: this.route });
  }
}