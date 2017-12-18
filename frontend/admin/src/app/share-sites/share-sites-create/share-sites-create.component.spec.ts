import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShareSitesCreateComponent } from './share-sites-create.component';

describe('ShareSitesCreateComponent', () => {
  let component: ShareSitesCreateComponent;
  let fixture: ComponentFixture<ShareSitesCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShareSitesCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShareSitesCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
