import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppComponent } from './app.component';
import { ApiService, UserService, CommonService, ConnectionService } from '../providers';
import { AppRouterRoutingModule } from '../router/app-router-routing.module';
import { FroalaEditorModule, FroalaViewModule } from 'angular2-froala-wysiwyg';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NoopInterceptor } from '../tools/http.Interceptor';

import {
  BoardsListComponent,
  ThreadsComponent,
  ThreadPageComponent,
  UserlistComponent,
  UserComponent,
  ConnectionComponent,
  AlertComponent,
  FixedButtonComponent,
  FabComponent,
  FabButton,
  ToTopComponent,
  ChipComponent,
  VoteBoxComponent
} from '../components';
import { SafeHTMLPipe, OrderByPipe, RepalcePipe } from '../pipes';
import { ClipDirective, IscrollDirective } from '../directives/index';
import { PopupModule } from '../providers/popup/popup.module';
import { LoadingModule } from '../providers/loading/loading.module';
import { AlertModule } from '../providers/alert/alert.module';
import { DialogModule } from '../providers/dialog/dialog.module';

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    HttpModule,
    FormsModule,
    ReactiveFormsModule,
    AppRouterRoutingModule,
    NgbModule.forRoot(),
    FroalaEditorModule.forRoot(),
    FroalaViewModule.forRoot(),
    PopupModule.forRoot(),
    LoadingModule.forRoot(),
    AlertModule.forRoot(),
    DialogModule.forRoot()
  ],
  declarations: [
    AppComponent,
    BoardsListComponent,
    ThreadsComponent,
    ThreadPageComponent,
    UserlistComponent,
    UserComponent,
    ConnectionComponent,
    AlertComponent,
    FixedButtonComponent,
    FabComponent,
    FabButton,
    FixedButtonComponent,
    ToTopComponent,
    ChipComponent,
    VoteBoxComponent,
    // Pipes
    SafeHTMLPipe,
    OrderByPipe,
    RepalcePipe,
    // Directives
    ClipDirective,
    IscrollDirective,
  ],
  entryComponents: [AlertComponent, FixedButtonComponent, ToTopComponent],
  providers: [CommonService, ApiService, UserService, ConnectionService, {
    provide: HTTP_INTERCEPTORS,
    useClass: NoopInterceptor,
    multi: true,
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
