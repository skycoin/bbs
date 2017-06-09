import { NgModule } from '@angular/core';
import { Routes, RouterModule, RouterOutletMap } from '@angular/router';
import { BoardsListComponent, ThreadsComponent, AddComponent, ThreadPageComponent } from "../components";

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
  { path: '**', redirectTo: '' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  // providers: [RouterOutletMap]
})
export class AppRouterRoutingModule { }
