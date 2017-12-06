import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
// import { HttpModule } from '@angular/http';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { ApiService, CommonService, UserService } from '../providers';
import { AppRouterRoutingModule } from '../router/app-router-routing.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NoopInterceptor } from '../tools/http.Interceptor';
import {
  BoardsListComponent,
  ThreadsComponent,
  ThreadPageComponent,
  ConnectionComponent
} from '../pages';
import {
  UserlistComponent,
  AlertComponent,
  FixedButtonComponent,
  FabComponent,
  FabButton,
  ToTopComponent,
  ChipComponent,
  SelectListComponent,
} from '../components';
import { SafeHTMLPipe, OrderByPipe, RepalcePipe, CutTextPipe } from '../pipes';
import { ClipDirective, EditorDirective, FocusDirective } from '../directives/index';
import { PopupModule } from '../providers/popup/popup.module';
import { LoadingModule } from '../providers/loading/loading.module';
import { AlertModule } from '../providers/alert/alert.module';
import { DialogModule } from '../providers/dialog/dialog.module';

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    // HttpModule,
    FormsModule,
    ReactiveFormsModule,
    AppRouterRoutingModule,
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
    ConnectionComponent,
    AlertComponent,
    FixedButtonComponent,
    FabComponent,
    FabButton,
    FixedButtonComponent,
    ToTopComponent,
    ChipComponent,
    SelectListComponent,
    // Pipes
    SafeHTMLPipe,
    OrderByPipe,
    RepalcePipe,
    CutTextPipe,
    // Directives
    ClipDirective,
    EditorDirective,
    FocusDirective
  ],
  entryComponents: [AlertComponent, FixedButtonComponent, ToTopComponent],
  providers: [CommonService, ApiService, UserService, {
    provide: HTTP_INTERCEPTORS,
    useClass: NoopInterceptor,
    multi: true,
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
