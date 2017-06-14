import { Component, OnInit, HostListener, HostBinding, ViewEncapsulation } from '@angular/core';
import { ApiService, ThreadPage, CommonService } from "../../providers";
import { Router, ActivatedRoute } from "@angular/router";
import { FormControl, FormGroup } from '@angular/forms';
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";
import { slideInLeftAnimation } from "../../animations/router.animations";

@Component({
  selector: 'threadPage',
  templateUrl: 'threadPage.html',
  styleUrls: ['threadPage.css'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation]
})

export class ThreadPageComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  private boardKey: string = '';
  private threadKey: string = '';
  private data: ThreadPage = { posts: [], thread: { name: '', description: '' } };
  private postForm = new FormGroup({
    title: new FormControl(),
    body: new FormControl()
  });
  private editorOptions = {
    placeholderText: 'Edit Your Content Here!',
    // toolbarButtons: ['bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', '|', 'fontFamily', 'fontSize', 'color', 'inlineStyle', 'paragraphStyle', '|', 'paragraphFormat', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '-', 'insertLink', 'insertImage', 'insertVideo', 'insertFile', 'insertTable', '|', 'emoticons', 'specialCharacters', 'insertHR', 'selectAll', 'clearFormatting', '|', 'print', 'spellChecker', 'help', 'html', '|', 'undo', 'redo'],
    toolbarButtons: ['bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', '|', 'fontFamily', 'fontSize', 'color', 'inlineStyle', 'paragraphStyle', '|', 'paragraphFormat', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '-', 'insertLink', '|', 'emoticons', 'specialCharacters', 'insertHR', 'selectAll', 'clearFormatting', '|', 'print', 'spellChecker', 'help', 'html', '|', 'undo', 'redo'],
    heightMin: 200,
    events: {
    },
  }
  constructor(
    private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    public common: CommonService) { }

  ngOnInit() {
    this.route.params.subscribe(res => {
      this.boardKey = res['board'];
      this.threadKey = res['thread'];
      this.open(this.boardKey, this.threadKey);
    })
  }

  openReply(content) {
    this.postForm.reset();
    this.modal.open(content, { backdrop: 'static', size: 'lg' }).result.then((result) => {
      let data = new FormData();
      if (result) {
        data.append('board', this.boardKey);
        data.append('thread', this.threadKey);
        data.append('title', this.postForm.get('title').value);
        data.append('body', this.postForm.get('body').value);
        this.api.addPost(data).subscribe(post => {
          if (post) {
            console.log('add post successfully:', post);
            this.data.posts.unshift(post);
            this.common.showAlert('Added successfully', 'success', 3000);
          }
        })
      }
    }, err => { });

  }

  reply() {
    if (!this.boardKey || !this.threadKey) {
      alert('Will not be able to post');
      return;
    }
    this.router.navigate(['/add', { exec: 'post', board: this.boardKey, thread: this.threadKey }]);
  }
  open(master, ref: string) {
    let data = new FormData();
    data.append('board', master);
    data.append('thread', ref);
    this.api.getThreadpage(data).subscribe(data => {
      this.data = data;
    });
  }

  @HostListener('window:scroll', ['$event'])
  windowScroll(event) {
    this.common.showOrHideToTopBtn();
  }

}