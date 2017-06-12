import { NgModule } from '@angular/core';
import { Routes, RouterModule, RouterOutletMap } from '@angular/router';
import {
  BoardsListComponent,
  ThreadsComponent,
  AddComponent,
  ThreadPageComponent,
  UserlistComponent,
  UserComponent,
  ConnectionComponent
} from "../components";

const routes: Routes = [
  { path: '', component: BoardsListComponent },
  {
    path: 'threads', children: [
      { path: '', component: ThreadsComponent },
      { path: 'p', component: ThreadPageComponent },
    ]
  },
  // { path: 'threads', component: ThreadsComponent },

  { path: 'add', component: AddComponent },
  { path: 'userlist', component: UserlistComponent },
  { path: 'user', component: UserComponent },
  { path: 'conn', component: ConnectionComponent },
  { path: '**', redirectTo: '' }

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  // providers: [RouterOutletMap]
})
export class AppRouterRoutingModule { }
