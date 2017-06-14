import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppComponent } from './app.component';
import { ApiService, UserService, CommonService, ConnectionService } from "../providers";
import { AppRouterRoutingModule } from "../router/app-router-routing.module";
import { FroalaEditorModule, FroalaViewModule } from 'angular2-froala-wysiwyg';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import {
  BoardsListComponent,
  ThreadsComponent,
  ThreadPageComponent,
  AddComponent,
  UserlistComponent,
  UserComponent,
  ConnectionComponent,
  AlertComponent,
} from "../components";
import { SafeHTMLPipe } from '../pipes';

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpModule,
    FormsModule,
    ReactiveFormsModule,
    AppRouterRoutingModule,
    NgbModule.forRoot(),
    FroalaEditorModule.forRoot(),
    FroalaViewModule.forRoot()
  ],
  declarations: [
    AppComponent,
    BoardsListComponent,
    ThreadsComponent,
    ThreadPageComponent,
    AddComponent,
    UserlistComponent,
    UserComponent,
    ConnectionComponent,
    AlertComponent,
    SafeHTMLPipe,
  ],
  entryComponents: [AlertComponent],
  providers: [CommonService, ApiService, UserService, ConnectionService],
  bootstrap: [AppComponent]
})
export class AppModule { }
