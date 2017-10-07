import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CollectionCreateComponent } from './collection-create.component';

describe('CollectionCreateComponent', () => {
  let component: CollectionCreateComponent;
  let fixture: ComponentFixture<CollectionCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CollectionCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CollectionCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
