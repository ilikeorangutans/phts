import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './collections/dashboard/dashboard.component';
import { DetailsComponent } from './collections/details/details.component';
import { PhotoDetailsComponent } from "./photos/photo-details/photo-details.component";

const routes: Routes = [
    {
        path: '',
        redirectTo: '/dashboard',
        pathMatch: 'full'
    },
    {
        path: 'dashboard',
        component: DashboardComponent
    },
    {
        path: 'collections/:slug',
        component: DetailsComponent
    },
    {
        path: 'collections/:slug/photos/:photoID',
        component: PhotoDetailsComponent
    }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutesModule {}
