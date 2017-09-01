import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from "@angular/router";
import { FormsModule } from "@angular/forms";
import { CollectionsDashboardComponent } from './collections-dashboard/collections-dashboard.component';
import { CollectionsFormComponent } from "./collections-form/collections-form.component";

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    FormsModule
  ],
  declarations: [
    CollectionsDashboardComponent,
    CollectionsFormComponent
  ],
  exports: [
    CollectionsDashboardComponent
  ]
})
export class CollectionsModule { }
