import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CollectionsFormComponent } from './collections-form.component';

describe('CollectionsFormComponent', () => {
  let component: CollectionsFormComponent;
  let fixture: ComponentFixture<CollectionsFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CollectionsFormComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CollectionsFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
