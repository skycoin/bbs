import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {
  BoardsListComponent,
  ThreadsComponent,
  ThreadPageComponent,
  UserlistComponent,
  UserComponent,
  ConnectionComponent
} from '../components';

const routes: Routes = [
  { path: '', component: BoardsListComponent, pathMatch: 'full' },
  {
    path: 'threads', children: [
      { path: '', component: ThreadsComponent },
      { path: 'p', component: ThreadPageComponent },
    ]
  },
  { path: 'userlist', component: UserlistComponent },
  { path: 'user', component: UserComponent },
  { path: 'conn', component: ConnectionComponent },
  { path: '**', redirectTo: '' }

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRouterRoutingModule { }
