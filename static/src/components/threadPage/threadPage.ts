import { Component, HostBinding, HostListener, OnInit, ViewEncapsulation } from '@angular/core';
import { ApiService, CommonService, ThreadPage, Post, VotesSummary, Thread } from '../../providers';
import { ActivatedRoute, Router } from '@angular/router';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
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
  public sort = 'esc';
  public boardKey = '';
  public threadKey = '';
  public data: ThreadPage = { posts: [], thread: { name: '', description: '' } };
  public postForm = new FormGroup({
    title: new FormControl('', Validators.required),
    body: new FormControl('', Validators.required),
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
    events: {},
  };

  constructor(private api: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private modal: NgbModal,
    private common: CommonService) {
  }

  ngOnInit() {
    this.route.params.subscribe(res => {
      this.boardKey = res['board'];
      this.threadKey = res['thread'];
      this.open(this.boardKey, this.threadKey);
    });
  }

  public setSort() {
    this.sort = this.sort === 'desc' ? 'esc' : 'desc';
  }
  addThreadVote(mode: string, thread: Thread, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (thread.uiOptions !== undefined && thread.uiOptions.voted !== undefined && thread.uiOptions.voted) {
      return;
    }
    thread.uiOptions = { voted: true };
    const data = new FormData();
    data.append('board', this.boardKey);
    data.append('thread', thread.ref);
    data.append('mode', mode);
    this.api.addThreadVote(data).subscribe(result => {
      if (result) {
        data.delete('mode');
        this.api.getThreadVotes(data).subscribe((votes: VotesSummary) => {
          thread.votes.up_votes = votes.up_votes;
          thread.votes.down_votes = votes.down_votes;
        }, err => {
          console.log('update vote fail');
        })
      } else {
        this.common.showErrorAlert('Vote Fail');
      }
    }, err => {
      thread.uiOptions.voted = false;
    })
  }
  addPostVote(mode: string, post: Post, ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (post.uiOptions !== undefined && post.uiOptions.voted !== undefined && post.uiOptions.voted) {
      return;
    }
    post.uiOptions = { voted: true };
    let data = new FormData();
    data.append('board', this.boardKey);
    data.append('post', post.ref);
    data.append('mode', mode);
    this.api.addPostVote(data).subscribe(result => {
      if (result) {
        data = new FormData();
        data.append('board', this.boardKey);
        data.append('post', post.ref);
        this.api.getPostVotes(data).subscribe((votes: VotesSummary) => {
          post.votes.up_votes = votes.up_votes;
          post.votes.down_votes = votes.down_votes;
        })
      } else {
        this.common.showErrorAlert('Vote Fail');
      }
    }, err => {
      post.uiOptions.voted = false;
    })
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
        });
      }
    }, err => {
    });

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
    }, err => {
      this.router.navigate(['']);
    });
  }

  @HostListener('window:scroll', ['$event'])
  windowScroll(event) {
    this.common.showOrHideToTopBtn(50);
  }
}
