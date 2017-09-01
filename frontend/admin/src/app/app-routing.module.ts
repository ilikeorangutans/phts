import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";

import { HomeComponent } from "./home/home.component";
import { CollectionsDashboardComponent } from "./collections/collections-dashboard/collections-dashboard.component";
import { CollectionsFormComponent } from "./collections/collections-form/collections-form.component";

const routes: Routes = [
    {
        path: '',
        redirectTo: '/home',
        pathMatch: 'full'
    },
    {
        path: 'home',
        component: HomeComponent
    },
    {
        path: 'collections',
        component: CollectionsDashboardComponent
    },
    {
        path: 'collections/new',
        component: CollectionsFormComponent
    }
];

@NgModule({
    imports: [ RouterModule.forRoot(routes) ],
    exports: [ RouterModule ]
})
export class AppRoutingModule {}