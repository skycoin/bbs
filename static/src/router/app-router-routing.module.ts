import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {
  BoardsListComponent,
  ConnectionComponent,
  ThreadPageComponent,
  ThreadsComponent,
} from '../pages';

const routes: Routes = [
  { path: '', component: BoardsListComponent, pathMatch: 'full' },
  {
    path: 'threads', children: [
      { path: '', component: ThreadsComponent },
      { path: 'p', component: ThreadPageComponent },
    ],
  },
  { path: 'conn', component: ConnectionComponent },
  { path: '**', redirectTo: '' },

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRouterRoutingModule {
}
