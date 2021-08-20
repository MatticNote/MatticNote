import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./page/index/index.module').then(m => m.IndexModule),
  },
  {
    path: '**',
    loadChildren: () => import('./page/not-found/not-found.module').then(m => m.NotFoundModule),
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
