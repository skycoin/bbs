import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { MaterialModule } from '@angular/material';

import { AppComponent } from './app.component';
import { ApiService, UserService, CommonService } from "../providers";
import { AppRouterRoutingModule } from "../router/app-router-routing.module";

import { BoardsListComponent, ThreadsComponent, ThreadPageComponent, AddComponent, UserlistComponent } from "../components";

@NgModule({
  declarations: [
    AppComponent, BoardsListComponent, ThreadsComponent, ThreadPageComponent, AddComponent, UserlistComponent
  ],
  imports: [
    BrowserModule, HttpModule, FormsModule, AppRouterRoutingModule
  ],
  providers: [CommonService, ApiService, UserService],
  bootstrap: [AppComponent]
})
export class AppModule { }
