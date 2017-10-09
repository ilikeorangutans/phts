import { ConfigurationListComponent } from './../renditions/configuration-list/configuration-list.component';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { CollectionBrowserComponent } from './../collections/collection-browser/collection-browser.component';
import { DashboardComponent } from '../dashboard/dashboard.component';
import { CollectionComponent } from '../collections/collection/collection.component';
import { CollectionDashboardComponent } from '../collections/collection-dashboard/collection-dashboard.component';
import { CollectionCreateComponent } from '../collections/collection-create/collection-create.component';
import { PhotoDetailsComponent } from '../photos/photo-details/photo-details.component';

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
        path: 'collections',
        component: CollectionBrowserComponent
    },
    {
        path: 'collections/new',
        component: CollectionCreateComponent
    },
    {
        path: 'collections/:slug',
        component: CollectionComponent,
        children: [
            {
                path: '',
                component: CollectionDashboardComponent
            },
            {
                path: 'photos/:photoID',
                component: PhotoDetailsComponent
            },
            {
                path: 'rendition-configurations',
                component: ConfigurationListComponent
            }
        ]
    }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutesModule {}
