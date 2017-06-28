import { Component, OnInit, HostListener, HostBinding, ViewEncapsulation } from '@angular/core';
import { ApiService, ThreadPage, CommonService } from '../../providers';
import { Router, ActivatedRoute } from '@angular/router';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModal, ModalDismissReasons } from '@ng-bootstrap/ng-bootstrap';
import { slideInLeftAnimation } from '../../animations/router.animations';

@Component({
  selector: 'app-threadpage',
  templateUrl: 'threadPage.html',
  styleUrls: ['threadPage.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation],
})

export class ThreadPageComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  public sort = 'desc';
  public boardKey = '';
  public threadKey = '';
  public data: ThreadPage = { posts: [], thread: { name: '', description: '' } };
  public postForm = new FormGroup({
    title: new FormControl('', Validators.required),
    body: new FormControl('', Validators.required)
  });
  public editorOptions = {
    placeholderText: 'Edit Your Content Here!',
    quickInsertButtons: ['table', 'ul', 'ol', 'hr'],
    toolbarButtons: [
      'bold',
      'italic',
      'underline',
      'strikeThrough',
      'subscript',
      'superscript',
      '|',
      'fontFamily',
      'fontSize',
      'color',
      'inlineStyle',
      'paragraphStyle',
      '|',
      'paragraphFormat',
      'align',
      'formatOL',
      'formatUL',
      'outdent',
      'indent',
      'quote',
      '-',
      'emoticons',
      'insertLink',
      '|',
      'insertHR',
      'selectAll',
      'clearFormatting',
      '|',
      'print',
      'spellChecker',
      'help',
      'html',
      '|',
      'undo',
      'redo'],
    heightMin: 200,
    events: {
    },
  }
  constructor(
    private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    private common: CommonService) { }

  ngOnInit() {
    this.route.params.subscribe(res => {
      this.boardKey = res['board'];
      this.threadKey = res['thread'];
      this.open(this.boardKey, this.threadKey);
    })
  }
  public setSort() {
    this.sort = this.sort === 'desc' ? 'esc' : 'desc';
  }
  openReply(content) {
    this.postForm.reset();
    this.modal.open(content, { backdrop: 'static', size: 'lg', keyboard: false }).result.then((result) => {
      if (result) {
        if (!this.postForm.valid) {
          this.common.showErrorAlert('Can not reply,title and content can not be empty');
          return;
        }
        const data = new FormData();
        data.append('board', this.boardKey);
        data.append('thread', this.threadKey);
        data.append('title', this.postForm.get('title').value);
        data.append('body', this.postForm.get('body').value);
        this.common.loading.start();
        this.api.addPost(data).subscribe(post => {
          if (post) {
            if (this.data.posts.length > 0) {
              this.data.posts.unshift(post);
            } else {
              this.data.posts = this.data.posts.concat(post);
            }
            this.common.loading.close();
            this.common.showAlert('Added successfully', 'success', 3000);
          }
        })
      }
    }, err => { });

  }
  open(master, ref: string) {
    if (master === '' || ref === '') {
      this.common.showErrorAlert('Parameter error!!!');
      return;
    }
    this.common.loading.start();
    const data = new FormData();
    data.append('board', master);
    data.append('thread', ref);
    this.api.getThreadpage(data).subscribe(res => {
      this.data = res;
      this.common.loading.close();
    });
  }

  @HostListener('window:scroll', ['$event'])
  windowScroll(event) {
    this.common.showOrHideToTopBtn(50);
  }
}
