import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppComponent } from './app.component';
import { ApiService, UserService, CommonService, ConnectionService } from "../providers";
import { AppRouterRoutingModule } from "../router/app-router-routing.module";

import {
  BoardsListComponent,
  ThreadsComponent,
  ThreadPageComponent,
  AddComponent,
  UserlistComponent,
  UserComponent,
  ModalComponent,
  ConnectionComponent,
  AlertComponent,
} from "../components";

@NgModule({
  imports: [
    BrowserModule,
    HttpModule,
    FormsModule,
    ReactiveFormsModule,
    AppRouterRoutingModule,
    NgbModule.forRoot(),
  ],
  declarations: [
    AppComponent,
    BoardsListComponent,
    ThreadsComponent,
    ThreadPageComponent,
    AddComponent,
    UserlistComponent,
    UserComponent,
    ModalComponent,
    ConnectionComponent,
    AlertComponent,
  ],
  entryComponents: [ModalComponent, AlertComponent],
  providers: [CommonService, ApiService, UserService, ConnectionService],
  bootstrap: [AppComponent]
})
export class AppModule { }
