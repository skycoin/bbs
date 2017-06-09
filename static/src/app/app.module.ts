import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { ApiService } from "../providers";
import { AppRouterRoutingModule } from "../router/app-router-routing.module";

import { BoardsListComponent, ThreadsComponent, ThreadPageComponent, AddComponent } from "../components";
@NgModule({
  declarations: [
    AppComponent, BoardsListComponent, ThreadsComponent, ThreadPageComponent, AddComponent
  ],
  imports: [
    BrowserModule, HttpModule, FormsModule,AppRouterRoutingModule
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule { }
