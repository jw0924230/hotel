<template>
  <div class="admin-page" :class="{ 'dark-mode': isDarkMode }">
    <header class="admin-header">
      <div class="header-container">
        <h1 class="admin-title">全台旅館管理後台</h1>
        <div class="header-actions">
          <span v-if="token" class="user-badge">
            <svg
              class="user-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
              />
            </svg>
            {{ loginUserEmail }} ({{
              userRole === "admin" ? "系統管理員" : "廠商帳號"
            }})
          </span>

          <button v-if="token" class="btn-logout" @click="logout">登出</button>
          <NuxtLink to="/" class="btn-home" v-if="userRole !== 'admin'"
            >回前台首頁</NuxtLink
          >
        </div>
      </div>
    </header>

    <!-- Login Form (if not authenticated) -->
    <div v-if="!token" class="login-wrapper">
      <div class="login-card">
        <div class="login-header">
          <span class="login-icon">
            <svg
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </span>
          <h2>後台管理登入</h2>
          <p>請輸入管理員或廠商帳號密碼以繼續</p>
        </div>
        <form @submit.prevent="login">
          <div v-if="loginError" class="alert-danger">{{ loginError }}</div>
          <div class="form-group">
            <label>電子信箱</label>
            <input
              v-model="loginEmail"
              type="email"
              required
              placeholder="example@gmail.com"
            />
          </div>
          <div class="form-group">
            <label>密碼</label>
            <input
              v-model="loginPassword"
              type="password"
              required
              placeholder="請輸入密碼"
            />
          </div>
          <button type="submit" class="btn-login" :disabled="loggingIn">
            {{ loggingIn ? "登入中..." : "登入後台" }}
          </button>
        </form>
        <div class="login-footer">
          <small>預設帳密請參閱系統設定檔</small>
        </div>
      </div>
    </div>

    <!-- Admin Workspace (if authenticated) -->
    <div v-else class="admin-main-wrapper">
      <!-- Side Navigation (Visible only to admin role for dashboard tab switching) -->
      <aside class="admin-nav" v-if="userRole === 'admin'">
        <div class="nav-brand">CMS 控制台</div>
        <nav class="nav-links">
          <button
            :class="['nav-link', { active: activeSection === 'hotels' }]"
            @click="activeSection = 'hotels'"
          >
            <svg
              class="nav-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
            <span>旅館資料管理</span>
          </button>
          <button
            :class="['nav-link', { active: activeSection === 'posts' }]"
            @click="
              activeSection = 'posts';
              fetchPosts();
            "
          >
            <svg
              class="nav-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
            <span>文章資料管理</span>
          </button>
          <button
            :class="['nav-link', { active: activeSection === 'homepage' }]"
            @click="
              activeSection = 'homepage';
              fetchHomepageHotels();
            "
          >
            <svg
              class="nav-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
              />
            </svg>
            <span>首頁版位管理</span>
          </button>
          <button
            :class="['nav-link', { active: activeSection === 'users' }]"
            @click="
              activeSection = 'users';
              fetchUsers();
            "
          >
            <svg
              class="nav-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
            <span>後台帳號管理</span>
          </button>
        </nav>
        <div class="nav-footer">
          <button
            v-if="token && userRole === 'admin'"
            class="btn-deploy"
            :disabled="isDeploying"
            @click="openDeployModal"
          >
            <svg
              class="deploy-btn-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
              style="width: 16px; height: 16px;"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
              />
            </svg>
            <span>部署前端網站</span>
          </button>
          <button
            class="btn-theme-toggle"
            @click="toggleDarkMode"
            aria-label="切換顯示模式"
          >
            <svg
              v-if="isDarkMode"
              class="theme-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m0-12.728l.707.707m12.728 12.728l.707-.707M12 8a4 4 0 100 8 4 4 0 000-8z"
              />
            </svg>
            <svg
              v-else
              class="theme-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
              />
            </svg>
            <span>{{ isDarkMode ? "淺色模式" : "深色模式" }}</span>
          </button>
          <NuxtLink to="/" class="btn-home-nav">回前台首頁</NuxtLink>
        </div>
      </aside>

      <!-- Workspace Container -->
      <div class="admin-workspace">
        <!-- ==================== 1. Hotels Workspace ==================== -->
        <div v-if="activeSection === 'hotels'" class="hotels-workspace">
          <!-- Sidebar / Hotel List -->
          <aside class="admin-sidebar">
            <div class="search-box">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="搜尋名稱或地址..."
                @input="debounceSearch"
                class="input-search"
              />
              <!-- Two-Panel Region/City Selector -->
              <div class="region-selector-wrapper">
                <button
                  class="btn-region-trigger"
                  @click="showRegionPicker = !showRegionPicker"
                >
                  <svg
                    class="filter-icon"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
                    />
                  </svg>
                  <span v-if="selectedCities.length === 0">全部地區</span>
                  <span v-else>{{ selectedCities.join(", ") }}</span>
                  <svg
                    class="chevron-icon"
                    :class="{ open: showRegionPicker }"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </button>
                <div v-if="showRegionPicker" class="region-picker-popover">
                  <div class="region-picker-body">
                    <div class="region-list">
                      <button
                        v-for="(region, idx) in regions"
                        :key="idx"
                        :class="[
                          'region-item',
                          { active: activeRegionIdx === idx },
                        ]"
                        @click="activeRegionIdx = idx"
                      >
                        {{ region.name }}
                      </button>
                    </div>
                    <div class="city-list" v-if="regions.length > 0">
                      <label
                        v-for="city in regions[activeRegionIdx]?.cities || []"
                        :key="city"
                        class="city-checkbox"
                      >
                        <input
                          type="checkbox"
                          :value="city"
                          v-model="selectedCities"
                          @change="onCityFilterChange"
                        />
                        <span>{{ city }}</span>
                      </label>
                    </div>
                  </div>
                  <div class="region-picker-footer">
                    <button class="btn-clear-filter" @click="clearCityFilter">
                      <svg
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2"
                        style="width: 14px; height: 14px"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                        />
                      </svg>
                      清除條件
                    </button>
                    <button
                      class="btn-close-picker"
                      @click="showRegionPicker = false"
                    >
                      關閉
                    </button>
                  </div>
                </div>
              </div>
              <button
                type="button"
                class="btn-create-hotel"
                @click="initNewHotel"
              >
                <svg
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                  style="width: 16px; height: 16px"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 4v16m8-8H4"
                  />
                </svg>
                新增旅館資料
              </button>
            </div>

            <div v-if="loadingList" class="loading-state">載入中...</div>
            <div v-else-if="hotels.length === 0" class="empty-state">
              無符合的旅館
            </div>
            <ul v-else class="hotel-list">
              <li
                v-for="h in hotels"
                :key="h.id"
                :class="['hotel-item', { active: selectedHotel?.id === h.id }]"
                @click="selectHotel(h.id)"
              >
                <div class="hotel-item-name">{{ h.name }}</div>
                <div class="hotel-item-info">
                  <span class="badge">{{ h.area || "未設定地區" }}</span>
                </div>
              </li>
            </ul>

            <!-- Pagination -->
            <div class="pagination" v-if="totalPages > 1">
              <button
                :disabled="currentPage === 1"
                @click="changePage(currentPage - 1)"
              >
                上一頁
              </button>
              <span>{{ currentPage }} / {{ totalPages }}</span>
              <button
                :disabled="currentPage === totalPages"
                @click="changePage(currentPage + 1)"
              >
                下一頁
              </button>
            </div>
          </aside>

          <!-- Main Editor -->
          <main class="admin-main">
            <div v-if="!selectedHotel" class="no-selection">
              <div class="icon">
                <svg
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
                  />
                </svg>
              </div>
              <h2>請從左側選擇要編輯的旅館</h2>
              <p>您也可以在下方手動輸入 ID 來編輯特定旅館：</p>
              <div class="quick-load">
                <input
                  v-model="quickId"
                  type="text"
                  placeholder="輸入旅館 ID，例如 9 或 1264"
                />
                <button @click="selectHotel(quickId)">載入</button>
              </div>
            </div>

            <div v-else class="editor-wrapper">
              <div class="editor-header">
                <h2>
                  {{ editForm.name }}
                  <span class="id-tag">ID: {{ editForm.id }}</span>
                </h2>
                <div class="actions" style="display: flex; gap: 8px; align-items: center;">
                  <button
                    type="button"
                    class="btn-preview"
                    @click="openHotelPreviewModal"
                    style="margin: 0;"
                  >
                    預覽頁面
                  </button>
                  <button class="btn-save" :disabled="saving" @click="saveHotel" style="margin: 0;">
                    {{ saving ? "儲存中..." : "儲存變更" }}
                  </button>
                </div>

              <div
                class="hotel-editor-tabs"
                role="tablist"
                aria-label="旅館資料編輯分頁"
              >
                <button
                  v-for="tab in hotelEditorTabs"
                  :key="tab.id"
                  type="button"
                  role="tab"
                  :aria-selected="activeHotelTab === tab.id"
                  :class="{ active: activeHotelTab === tab.id }"
                  @click="setHotelEditorTab(tab.id)"
                >
                  <span>{{ tab.label }}</span>
                </button>
              </div>

              <!-- Editor vertical sections list -->
              <div class="editor-vertical-list">
                <!-- Card 1: Basic Info -->
                <div v-show="activeHotelTab === 'basic'" class="card">
                  <h3>基本資訊</h3>
                  <div
                    class="form-grid-2"
                    v-if="isNewHotel"
                    style="margin-bottom: 16px"
                  >
                    <div class="form-group">
                      <label>旅館 ID <span class="required">*</span></label>
                      <input
                        v-model="editForm.id"
                        type="text"
                        placeholder="請輸入旅館 ID，例如 1265"
                        required
                      />
                    </div>
                  </div>
                  <div class="form-grid-2">
                    <div class="form-group">
                      <label>旅館名稱 <span class="required">*</span></label>
                      <input
                        v-model="editForm.name"
                        type="text"
                        placeholder="例如: 丹迪旅店 大安森林公園店"
                        required
                      />
                    </div>
                    <div class="form-group">
                      <label>地區分類 <span class="required">*</span></label>
                      <div class="region-selector-wrapper">
                        <button
                          type="button"
                          class="btn-region-trigger"
                          @click="showEditRegionPicker = !showEditRegionPicker"
                        >
                          <svg
                            class="filter-icon"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
                            />
                          </svg>
                          <span>{{ editForm.area || "請選擇地區" }}</span>
                          <svg
                            class="chevron-icon"
                            :class="{ open: showEditRegionPicker }"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              d="M19 9l-7 7-7-7"
                            />
                          </svg>
                        </button>
                        <div
                          v-if="showEditRegionPicker"
                          class="region-picker-popover"
                        >
                          <div class="region-picker-body">
                            <div class="region-list">
                              <button
                                type="button"
                                v-for="(region, idx) in regions"
                                :key="idx"
                                :class="[
                                  'region-item',
                                  { active: activeEditRegionIdx === idx },
                                ]"
                                @click="activeEditRegionIdx = idx"
                              >
                                {{ region.name }}
                              </button>
                            </div>
                            <div class="city-list" v-if="regions.length > 0">
                              <button
                                type="button"
                                v-for="city in regions[activeEditRegionIdx]
                                  ?.cities || []"
                                :key="city"
                                :class="[
                                  'city-picker-item',
                                  { active: editForm.area === city },
                                ]"
                                @click="selectEditArea(city)"
                              >
                                {{ city }}
                              </button>
                            </div>
                          </div>
                          <div
                            class="region-picker-footer"
                            style="justify-content: flex-end"
                          >
                            <button
                              type="button"
                              class="btn-close-picker"
                              @click="showEditRegionPicker = false"
                            >
                              關閉
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="form-group">
                    <label>完整地址 <span class="required">*</span></label>
                    <input
                      v-model="editForm.address"
                      type="text"
                      placeholder="例如: 台北市大安區信義路三段33號"
                      required
                    />
                  </div>

                  <div class="form-grid-3">
                    <div class="form-group">
                      <label>聯絡電話</label>
                      <input
                        v-model="editForm.phone"
                        type="text"
                        placeholder="例如: 02-2707-6899"
                      />
                    </div>
                    <div class="form-group">
                      <label>傳真號碼</label>
                      <input
                        v-model="editForm.fax"
                        type="text"
                        placeholder="例如: 02-2707-6899"
                      />
                    </div>
                    <div class="form-group">
                      <label>官方網站</label>
                      <input
                        v-model="editForm.website"
                        type="text"
                        placeholder="例如: https://www.dandyhotel.com.tw/"
                      />
                    </div>
                  </div>

                  <div class="form-group">
                    <label>電子信箱</label>
                    <input
                      v-model="editForm.email"
                      type="text"
                      placeholder="例如: service@dandyhotel.com"
                    />
                  </div>

                  <div class="form-group" style="margin-top: 16px">
                    <label class="checkbox-label">
                      <input type="checkbox" v-model="editForm.is_disabled" />
                      <span
                        >停用此旅館 (勾選後前台 API
                        將排除此旅館，前台網頁將不顯示)</span
                      >
                    </label>
                  </div>
                </div>

                <!-- Card 2: Stay Info -->
                <div v-show="activeHotelTab === 'booking'" class="card">
                  <h3>訂房按鈕與連結</h3>
                  <div class="form-group">
                    <label>一鍵訂房/預定連結 (訂房平台跳轉聯播網 URL)</label>
                    <input
                      v-model="editForm.booking_link"
                      type="text"
                      placeholder="輸入 FunNow 或其他訂房網站的分潤預定連結..."
                    />
                  </div>
                </div>

                <!-- Card 3: Price Info -->
                <div v-show="activeHotelTab === 'booking'" class="card">
                  <h3>價格設定</h3>
                  <div class="pricing-settings">
                    <section class="pricing-section">
                      <div class="pricing-section-heading">
                        <span class="pricing-section-icon">住</span>
                        <div>
                          <h4>住宿價格</h4>
                          <p>設定平日與假日的每晚價格</p>
                        </div>
                      </div>
                      <div class="pricing-option-grid">
                        <div class="pricing-option">
                          <span class="pricing-option-label">平日</span>
                          <div class="form-group">
                            <label>住宿價格</label>
                            <IntegerInput
                              v-model="editForm.pricing.weekday_stay"
                            />
                          </div>
                        </div>
                        <div class="pricing-option">
                          <span class="pricing-option-label holiday">假日</span>
                          <div class="form-group">
                            <label>住宿價格</label>
                            <IntegerInput
                              v-model="editForm.pricing.holiday_stay"
                            />
                          </div>
                        </div>
                      </div>
                    </section>

                    <section class="pricing-section">
                      <div class="pricing-section-heading">
                        <span class="pricing-section-icon rest">休</span>
                        <div>
                          <h4>休息方案</h4>
                          <p>平日與假日分別設定時數及價格</p>
                        </div>
                      </div>
                      <div class="pricing-option-grid">
                        <div class="pricing-option">
                          <span class="pricing-option-label">平日</span>
                          <div class="pricing-field-pair">
                            <div class="form-group">
                              <label>休息時數</label>
                              <IntegerInput
                                v-model="editForm.pricing.weekday_rest_hours"
                              />
                            </div>
                            <div class="form-group">
                              <label>休息價格</label>
                              <IntegerInput
                                v-model="editForm.pricing.weekday_rest"
                              />
                            </div>
                          </div>
                        </div>
                        <div class="pricing-option">
                          <span class="pricing-option-label holiday">假日</span>
                          <div class="pricing-field-pair">
                            <div class="form-group">
                              <label>休息時數</label>
                              <IntegerInput
                                v-model="editForm.pricing.holiday_rest_hours"
                              />
                            </div>
                            <div class="form-group">
                              <label>休息價格</label>
                              <IntegerInput
                                v-model="editForm.pricing.holiday_rest"
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </section>
                  </div>
                </div>

                <!-- Card 4: Description (Rich Text Editor) -->
                <div v-show="activeHotelTab === 'content'" class="card">
                  <h3>簡介</h3>
                  <WangEditor
                    v-model="editForm.description"
                    placeholder="請輸入旅館簡介…"
                  />
                </div>

                <!-- Card 5: Rules (Rich Text Editor) -->
                <div v-show="activeHotelTab === 'content'" class="card">
                  <h3>住房須知與規定</h3>
                  <WangEditor
                    v-model="editForm.housing_rules"
                    placeholder="請輸入住房須知與規定…"
                  />
                </div>

                <!-- Card 7: Images Manager -->
                <div v-show="activeHotelTab === 'images'" class="card">
                  <h3>圖片管理 (Imgur 圖床連結，最多 10 張)</h3>
                  <div class="image-editor-section">
                    <div
                      v-if="!editForm.images || editForm.images.length === 0"
                      class="no-images"
                    >
                      目前無圖片，請點擊下方按鈕新增圖片網址。
                    </div>
                    <div class="image-urls-list" v-else>
                      <div
                        v-for="(img, idx) in editForm.images"
                        :key="idx"
                        :class="[
                          'image-url-row',
                          { dragging: draggedImageIndex === idx },
                        ]"
                        draggable="true"
                        @dragstart="handleImageDragStart($event, idx)"
                        @dragover.prevent="handleImageDragOver($event, idx)"
                        @drop="handleImageDrop($event, idx)"
                        @dragend="handleImageDragEnd"
                      >
                        <div class="image-url-header">
                          <span class="img-order-badge"
                            >圖片 {{ idx + 1 }}</span
                          >
                          <div class="row-actions">
                            <span
                              class="drag-handle"
                              title="按住拖曳排序"
                              style="
                                cursor: grab;
                                margin-right: 12px;
                                color: #94a3b8;
                                font-size: 16px;
                                user-select: none;
                                display: flex;
                                align-items: center;
                              "
                              >☰</span
                            >
                            <button
                              type="button"
                              class="btn-row-preview"
                              :disabled="!editForm.images[idx]"
                              @click="previewImage(editForm.images[idx])"
                            >
                              預覽
                            </button>
                            <button
                              type="button"
                              class="btn-row-delete"
                              @click="deleteImage(idx)"
                            >
                              刪除
                            </button>
                          </div>
                        </div>
                        <div
                          class="image-input-container"
                          style="
                            display: flex;
                            gap: 12px;
                            align-items: center;
                            margin-top: 8px;
                            width: 100%;
                          "
                        >
                          <input
                            v-model="editForm.images[idx]"
                            type="text"
                            placeholder="請輸入 imgur 圖片網址 (例如: https://i.imgur.com/xxxx.jpg)"
                            @blur="validateAndFormatImgur(idx)"
                            class="input-image-url"
                            style="flex: 1; margin: 0; min-width: 0"
                          />
                          <div
                            class="image-thumbnail-preview"
                            style="
                              width: 50px;
                              height: 50px;
                              border-radius: 6px;
                              overflow: hidden;
                              border: 1px solid #cbd5e1;
                              display: flex;
                              align-items: center;
                              justify-content: center;
                              background-color: #f1f5f9;
                              flex-shrink: 0;
                            "
                          >
                            <img
                              v-if="
                                editForm.images[idx] &&
                                editForm.images[idx].trim() !== ''
                              "
                              :src="getImgurThumbnail(editForm.images[idx])"
                              alt="縮圖"
                              style="
                                width: 100%;
                                height: 100%;
                                object-fit: cover;
                              "
                              @error="handleThumbnailError($event)"
                            />
                            <span v-else style="color: #94a3b8; font-size: 11px"
                              >無圖</span
                            >
                          </div>
                        </div>
                      </div>
                    </div>
                    <button
                      type="button"
                      class="btn-add-img"
                      v-if="!editForm.images || editForm.images.length < 10"
                      @click="addImageUrl"
                    >
                      新增圖片連結
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </main>
        </div>

        <!-- ==================== 2. Users Workspace (Admin Only) ==================== -->
        <div
          v-else-if="activeSection === 'users' && userRole === 'admin'"
          class="users-workspace"
        >
          <div class="workspace-header">
            <h2>後台帳號權限管理</h2>
            <button class="btn-create-user" @click="openCreateUserModal">
              新增後台帳號
            </button>
          </div>

          <div v-if="loadingUsers" class="loading-state">載入帳號資料中...</div>
          <div v-else class="users-table-wrapper">
            <table class="users-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>電子信箱</th>
                  <th>角色權限</th>
                  <th>建立時間</th>
                  <th>操作設定</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="u in users" :key="u.id">
                  <td>{{ u.id }}</td>
                  <td>
                    <strong>{{ u.email }}</strong>
                  </td>
                  <td>
                    <span :class="['role-badge', u.role]">
                      {{ u.role === "admin" ? "管理者" : "廠商" }}
                    </span>
                  </td>
                  <td>{{ formatDate(u.created_at) }}</td>
                  <td>
                    <div class="table-actions">
                      <button
                        class="btn-table-edit"
                        @click="openEditUserModal(u)"
                      >
                        編輯
                      </button>
                      <button
                        class="btn-table-delete"
                        @click="deleteUser(u.id)"
                      >
                        刪除
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- User Editor Modal Overlay -->
          <div class="modal-overlay" v-if="showUserModal">
            <div class="modal-content">
              <div class="modal-header">
                <h3>
                  {{ editingUser ? "編輯帳號設定" : "新增管理員/廠商帳號" }}
                </h3>
                <button class="btn-close-modal" @click="showUserModal = false">
                  ×
                </button>
              </div>
              <form @submit.prevent="saveUser">
                <div v-if="userFormError" class="alert-danger">
                  {{ userFormError }}
                </div>
                <div class="form-group">
                  <label
                    >電子信箱 (Email) <span class="required">*</span></label
                  >
                  <input
                    v-model="userFormEmail"
                    type="email"
                    required
                    placeholder="example@gmail.com"
                  />
                </div>
                <div class="form-group">
                  <label
                    >密碼 {{ editingUser ? "(留空表示不修改)" : ""
                    }}<span v-if="!editingUser" class="required">*</span></label
                  >
                  <input
                    v-model="userFormPassword"
                    type="password"
                    :required="!editingUser"
                    placeholder="請輸入密碼"
                  />
                </div>
                <div class="form-group">
                  <label>角色權限 <span class="required">*</span></label>
                  <select v-model="userFormRole" class="select-role">
                    <option value="vendor">廠商 (僅能編輯旅館資料)</option>
                    <option value="admin">管理者 (可管理帳號與資料)</option>
                  </select>
                </div>
                <div class="modal-footer">
                  <button
                    type="button"
                    class="btn-cancel"
                    @click="showUserModal = false"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    class="btn-submit"
                    :disabled="savingUser"
                  >
                    {{ savingUser ? "儲存中..." : "確定儲存" }}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>

        <!-- ==================== 3. Posts Workspace (Admin Only) ==================== -->
        <div
          v-else-if="activeSection === 'posts' && userRole === 'admin'"
          class="posts-workspace"
        >
          <!-- Sidebar / Post List -->
          <aside class="admin-sidebar">
            <div class="search-box">
              <input
                v-model="postSearchQuery"
                type="text"
                placeholder="搜尋文章標題..."
                @input="fetchPosts"
                class="input-search"
              />
              <button class="btn-create-post" @click="openCreatePost">
                新增文章
              </button>
            </div>

            <div v-if="loadingPostsList" class="loading-state">載入中...</div>
            <div v-else-if="posts.length === 0" class="empty-state">
              目前無文章
            </div>
            <ul v-else class="post-list">
              <li
                v-for="p in posts"
                :key="p.id"
                :class="['post-item', { active: selectedPost?.id === p.id }]"
                @click="selectPost(p.id)"
              >
                <div class="post-item-title">{{ p.title }}</div>
                <div class="post-item-info">
                  <span class="tag-badge" v-for="tag in p.tags" :key="tag">{{
                    tag
                  }}</span>
                </div>
              </li>
            </ul>

            <!-- Pagination -->
            <div class="pagination" v-if="postTotalPages > 1">
              <button
                :disabled="postCurrentPage === 1"
                @click="
                  postCurrentPage--;
                  fetchPosts();
                "
              >
                上一頁
              </button>
              <span>{{ postCurrentPage }} / {{ postTotalPages }}</span>
              <button
                :disabled="postCurrentPage === postTotalPages"
                @click="
                  postCurrentPage++;
                  fetchPosts();
                "
              >
                下一頁
              </button>
            </div>
          </aside>

          <!-- Main Editor -->
          <main class="admin-main">
            <div v-if="!selectedPost" class="no-selection">
              <div class="icon">
                <svg
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
                  />
                </svg>
              </div>
              <h2>請從左側選擇文章，或點擊「新增文章」開始編輯</h2>
            </div>

            <div v-else class="editor-wrapper">
              <div class="editor-header">
                <h2>
                  {{ postEditForm.id ? "編輯文章" : "新增文章" }}
                  <span class="id-tag" v-if="postEditForm.id"
                    >ID: {{ postEditForm.id }}</span
                  >
                </h2>
                <div class="actions">
                  <button
                    class="btn-delete-post"
                    v-if="postEditForm.id"
                    @click="deletePost(postEditForm.id)"
                  >
                    刪除文章
                  </button>
                  <button
                    type="button"
                    class="btn-preview"
                    @click="openPostPreviewModal"
                  >
                    預覽頁面
                  </button>
                  <button
                    class="btn-save"
                    :disabled="postSaving"
                    @click="savePost"
                  >
                    {{ postSaving ? "儲存中..." : "儲存變更" }}
                  </button>
                </div>
              </div>

              <div class="editor-vertical-list">
                <!-- Card 1: Basic Info -->
                <div class="card">
                  <h3>文章基本資訊</h3>
                  <div class="form-group">
                    <label>文章標題 <span class="required">*</span></label>
                    <input
                      v-model="postEditForm.title"
                      type="text"
                      placeholder="請輸入文章標題"
                      required
                    />
                  </div>
                  <div class="form-group">
                    <label>封面圖片連結 (URL)</label>
                    <input
                      v-model="postEditForm.image"
                      type="text"
                      placeholder="請輸入 Imgur 圖片連結，例如: https://i.imgur.com/xxxx.jpg"
                    />
                    <div class="url-preview" v-if="postEditForm.image">
                      <img
                        :src="postEditForm.image"
                        alt="Cover Preview"
                        @error="handleImgPreviewError"
                      />
                    </div>
                  </div>
                </div>

                <!-- Card 2: Categorization Tags -->
                <div class="card">
                  <h3>分類標籤 (最多三個)</h3>
                  <div class="tags-input-container">
                    <div
                      class="tags-chips"
                      v-if="postEditForm.tags && postEditForm.tags.length > 0"
                    >
                      <span
                        class="tag-chip"
                        v-for="(tag, idx) in postEditForm.tags"
                        :key="idx"
                      >
                        {{ tag }}
                        <button
                          type="button"
                          class="btn-remove-tag"
                          @click="removeTag(idx)"
                        >
                          ×
                        </button>
                      </span>
                    </div>
                    <div v-else class="help-text">
                      目前無標籤。請在下方輸入並點選新增。
                    </div>
                    <div
                      class="tag-input-row"
                      v-if="!postEditForm.tags || postEditForm.tags.length < 3"
                    >
                      <input
                        v-model="newTagText"
                        type="text"
                        placeholder="輸入標籤名稱，例如: 住宿推薦"
                        @keyup.enter="addTag"
                      />
                      <button type="button" class="btn-add-tag" @click="addTag">
                        新增標籤
                      </button>
                    </div>
                    <small class="help-text"
                      >標籤將作為文章的分類使用。最多僅能設定三個標籤。</small
                    >
                  </div>
                </div>

                <!-- Card 3: SEO TKD 設定 -->
                <div class="card">
                  <h3>SEO 設定 (TKD)</h3>
                  <div class="form-group">
                    <label>SEO Title (頁面標題)</label>
                    <input
                      v-model="postEditForm.seo_title"
                      type="text"
                      placeholder="請輸入網頁 Meta Title"
                    />
                  </div>
                  <div class="form-group">
                    <label>SEO Keywords (關鍵字，以半形逗號分隔)</label>
                    <input
                      v-model="postEditForm.seo_keywords"
                      type="text"
                      placeholder="例如: 飯店推薦, 汽車旅館休息, 大安區"
                    />
                  </div>
                  <div class="form-group">
                    <label>SEO Description (網頁描述)</label>
                    <textarea
                      v-model="postEditForm.seo_description"
                      rows="3"
                      placeholder="請輸入網頁 Meta Description 簡介"
                    ></textarea>
                  </div>
                </div>

                <!-- Card 4: Content Editor -->
                <div class="card">
                  <h3>
                    文章內文 (富文本編輯器) <span class="required">*</span>
                  </h3>
                  <WangEditor
                    v-model="postEditForm.content"
                    placeholder="請輸入文章內容…"
                    height="420px"
                  />
                </div>

                <!-- Card 5: Ad / Promotion Link -->
                <div class="card">
                  <h3>廣告與推廣 HTML</h3>
                  <div class="form-group">
                    <label>文章底部的推廣廣告 HTML</label>
                    <textarea
                      v-model="postEditForm.ad_link"
                      rows="3"
                      placeholder="例如: 推薦信用卡: <a href='...'>精選信用卡</a>"
                    ></textarea>
                  </div>
                </div>
              </div>
            </div>
          </main>

          <!-- Post Deletion Confirmation Modal Overlay -->
          <div class="modal-overlay" v-if="showPostDeleteConfirmModal">
            <div class="modal-content" style="max-width: 450px">
              <div class="modal-header">
                <h3>確認刪除文章</h3>
                <button
                  class="btn-close-modal"
                  @click="showPostDeleteConfirmModal = false"
                >
                  ×
                </button>
              </div>
              <div class="modal-body" style="padding: 20px 0">
                <p
                  style="margin-bottom: 16px; color: #ef4444; font-weight: 600"
                >
                  ⚠️ 注意：您確定要刪除這篇文章嗎？此操作將永久刪除且無法復原。
                </p>
                <div class="form-group">
                  <label style="display: block; margin-bottom: 6px"
                    >請輸入 「<strong style="color: #ef4444">確認刪除</strong>」
                    以繼續進行：</label
                  >
                  <input
                    v-model="postDeleteConfirmInput"
                    type="text"
                    placeholder="請在此輸入「確認刪除」"
                    style="margin-top: 8px; width: 100%"
                    @keyup.enter="executeDeletePost"
                  />
                </div>
              </div>
              <div class="modal-footer">
                <button
                  type="button"
                  class="btn-cancel"
                  @click="showPostDeleteConfirmModal = false"
                >
                  取消
                </button>
                <button
                  type="button"
                  class="btn-submit"
                  :disabled="postDeleteConfirmInput !== '確認刪除'"
                  @click="executeDeletePost"
                  style="background-color: #ef4444; color: white"
                >
                  確定刪除
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- ==================== 4. Homepage Selections Workspace (Admin Only) ==================== -->
        <div
          v-else-if="activeSection === 'homepage' && userRole === 'admin'"
          class="homepage-workspace"
        >
          <div class="workspace-header">
            <h2>首頁城市精選推薦編輯</h2>
            <button
              class="btn-save"
              :disabled="savingHomepage || selectedCount !== 6"
              @click="saveHomepageHotels"
            >
              {{ savingHomepage ? "儲存中..." : "儲存首頁設定" }}
            </button>
          </div>

          <div v-if="loadingHomepage" class="loading-state">
            載入設定資料中...
          </div>
          <div v-else class="homepage-editor-container">
            <!-- City Selection Tabs -->
            <div class="city-tabs">
              <button
                v-for="city in homepageCities"
                :key="city"
                type="button"
                :class="[
                  'city-tab-btn',
                  { active: activeHomepageCity === city },
                ]"
                @click="
                  activeHomepageCity = city;
                  homepageSearchQuery = '';
                "
              >
                {{ city }}
              </button>
            </div>

            <!-- Transfer Layout Container -->
            <div class="transfer-container">
              <!-- Left panel: Available hotels in this city -->
              <div class="transfer-panel available-panel">
                <div class="panel-header-custom">
                  <h3>{{ activeHomepageCity }} 可選擇旅館</h3>
                  <div class="search-wrapper-custom">
                    <input
                      v-model="homepageSearchQuery"
                      type="text"
                      placeholder="輸入 ID 或旅館名稱搜尋..."
                      class="input-search-inline"
                    />
                  </div>
                </div>
                <div class="panel-body scrollable-y">
                  <div
                    v-for="h in filteredAvailableHotels"
                    :key="h.id"
                    class="transfer-item"
                  >
                    <div class="item-info">
                      <span class="item-id">ID: {{ h.id }}</span>
                      <strong class="item-name">{{ h.name }}</strong>
                      <span class="item-addr">{{ h.address || "無地址" }}</span>
                    </div>
                    <button
                      type="button"
                      class="btn-select-item"
                      :disabled="isHotelSelected(h.id) || selectedCount >= 6"
                      @click="selectHotelForSlots(h.id)"
                    >
                      選擇
                    </button>
                  </div>
                  <div
                    v-if="filteredAvailableHotels.length === 0"
                    class="empty-state-text"
                  >
                    無相符的旅館資料
                  </div>
                </div>
              </div>

              <!-- Right panel: Selected hotels (Slots 1-6) -->
              <div class="transfer-panel selected-panel">
                <div class="panel-header-custom">
                  <h3>已選擇推薦旅館 (已選 {{ selectedCount }} / 6)</h3>
                  <span
                    :class="['status-badge', { valid: selectedCount === 6 }]"
                  >
                    {{
                      selectedCount === 6
                        ? "符合儲存條件"
                        : "須設定剛好 6 筆才能儲存"
                    }}
                  </span>
                </div>
                <div class="panel-body">
                  <div
                    v-for="index in 6"
                    :key="index"
                    :class="[
                      'slot-row',
                      { dragging: draggedIndex === index - 1 },
                    ]"
                    draggable="true"
                    @dragstart="handleDragStart($event, index - 1)"
                    @dragover.prevent="handleDragOver($event, index - 1)"
                    @drop="handleDrop($event, index - 1)"
                    @dragend="handleDragEnd"
                  >
                    <span class="slot-number-badge">版位 {{ index }}</span>
                    <div class="slot-content">
                      <div
                        v-if="getSelectedHotelInSlot(index - 1)"
                        class="selected-hotel-display"
                      >
                        <div class="display-text">
                          <span class="item-id"
                            >ID:
                            {{ getSelectedHotelInSlot(index - 1).id }}</span
                          >
                          <strong class="display-name">{{
                            getSelectedHotelInSlot(index - 1).name
                          }}</strong>
                        </div>
                        <div class="slot-actions">
                          <span
                            class="drag-handle"
                            title="按住拖曳排序"
                            style="
                              cursor: grab;
                              margin-right: 12px;
                              color: #94a3b8;
                              font-size: 16px;
                              user-select: none;
                              display: flex;
                              align-items: center;
                            "
                            >☰</span
                          >
                          <button
                            type="button"
                            class="btn-slot-remove"
                            @click="removeHotelFromSlot(index - 1)"
                          >
                            移除
                          </button>
                        </div>
                      </div>
                      <div v-else class="empty-slot-placeholder">
                        版位空缺 (請自左側選擇)
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- Global Floating Snackbar -->
      <div :class="['snackbar', snackbar.type, { show: snackbar.show }]">
        <div class="snackbar-content">
          <svg
            v-if="snackbar.type === 'success'"
            class="snackbar-icon"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <svg
            v-else
            class="snackbar-icon"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
          <span>{{ snackbar.message }}</span>
        </div>
        <button
          type="button"
          class="snackbar-close"
          @click="snackbar.show = false"
        >
          ×
        </button>
      </div>
    </div>

    <!-- Frontend Deployment Modal Overlay -->
    <div class="modal-overlay" v-if="showDeployConfirmModal">
      <div class="modal-content" style="max-width: 500px">
        <div class="modal-header">
          <h3>部署前端網站</h3>
          <button
            class="btn-close-modal"
            @click="closeDeployModal"
            :disabled="isDeploying"
          >
            ×
          </button>
        </div>
        <div class="modal-body" style="padding: 20px 0; text-align: center;">
          <div v-if="deployStatus === 'idle'">
            <div style="font-size: 40px; margin-bottom: 12px;">🚀</div>
            <p style="font-size: 16px; font-weight: 600; margin-bottom: 8px; color: #1e293b;">
              您確定要觸發前端網站部署嗎？
            </p>
            <p style="color: #64748b; font-size: 14px; line-height: 1.5; padding: 0 16px;">
              這將觸發 GitHub Actions。系統會從後端重新抓取最新資料並編譯產生靜態網站，整個過程通常需要 3~5 分鐘。
            </p>
          </div>
          <div v-else-if="deployStatus === 'deploying'">
            <div class="deploy-spinner" style="margin: 20px auto;"></div>
            <p style="font-size: 16px; font-weight: 600; color: #3b82f6;">
              正在發送部署請求...
            </p>
          </div>
          <div v-else-if="deployStatus === 'success'">
            <div style="font-size: 40px; margin-bottom: 12px; color: #10b981;">✓</div>
            <p style="font-size: 16px; font-weight: 600; color: #10b981; margin-bottom: 8px;">
              已成功觸發部署！
            </p>
            <p style="color: #64748b; font-size: 14px; line-height: 1.5; padding: 0 16px;">
              {{ deployMessage }}
            </p>
          </div>
          <div v-else-if="deployStatus === 'error'">
            <div style="font-size: 40px; margin-bottom: 12px; color: #ef4444;">⚠️</div>
            <p style="font-size: 16px; font-weight: 600; color: #ef4444; margin-bottom: 8px;">
              部署發送失敗
            </p>
            <p style="color: #ef4444; font-size: 14px; line-height: 1.5; padding: 0 16px;">
              {{ deployMessage }}
            </p>
          </div>
        </div>
        <div class="modal-footer" style="justify-content: center; gap: 12px;">
          <button
            type="button"
            class="btn-cancel"
            @click="closeDeployModal"
            v-if="deployStatus === 'idle' || deployStatus === 'success' || deployStatus === 'error'"
          >
            {{ deployStatus === 'idle' ? '取消' : '關閉' }}
          </button>
          <button
            type="button"
            class="btn-submit"
            @click="executeFrontendDeploy"
            v-if="deployStatus === 'idle'"
            style="background-color: #3b82f6; color: white;"
          >
            確定部署
          </button>
        </div>
      </div>
    </div>

    <!-- Post CSR Preview Modal Overlay -->
    <div class="modal-overlay" v-if="showPostPreviewModal" @click.self="showPostPreviewModal = false">
      <div class="modal-content" style="max-width: 850px; width: 95%; max-height: 85vh; display: flex; flex-direction: column; padding: 0;">
        <div class="modal-header" style="padding: 16px 24px; border-bottom: 1px solid #e2e8f0;">
          <h3 style="margin: 0; font-size: 18px;">文章前台模擬預覽 (CSR 畫面)</h3>
          <button class="btn-close-modal" @click="showPostPreviewModal = false">×</button>
        </div>
        <div class="modal-body" style="padding: 30px; overflow-y: auto; background: #f8f9fa; flex: 1;">
          <div style="max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 12px; box-shadow: 0 5px 20px rgba(0,0,0,0.05);">
            <div style="margin-bottom: 20px; color: #7f8c8d; font-size: 14px;">
              首頁 &gt; 部落格 &gt; <span style="color: #95a5a6;">{{ postEditForm.title || "未命名文章" }}</span>
            </div>
            
            <header style="margin-bottom: 30px; text-align: left; border-bottom: 1px solid #f1f5f9; padding-bottom: 20px;">
              <h1 style="font-size: 32px; color: #2C3E50; line-height: 1.4; margin: 0 0 15px 0; font-weight: 700;">
                {{ postEditForm.title || "請輸入文章標題" }}
              </h1>
              <div style="display: flex; justify-content: space-between; align-items: center;">
                <span style="background: #E74C3C; color: white; padding: 4px 12px; border-radius: 20px; font-size: 13px; font-weight: 600;">
                  {{ postEditForm.tags && postEditForm.tags.length > 0 ? postEditForm.tags[0] : '精選專欄' }}
                </span>
                <span style="color: #95A5A6; font-size: 14px;">
                  {{ new Date().toISOString().split('T')[0] }}
                </span>
              </div>
            </header>

            <div style="margin-bottom: 30px; border-radius: 8px; overflow: hidden;" v-if="postEditForm.image">
              <img :src="postEditForm.image" alt="Featured Image" style="width: 100%; height: auto; display: block;" />
            </div>

            <div class="article-preview-body" v-html="parsedPostContent" style="font-size: 18px; line-height: 1.8; color: #2c3e50;"></div>
            
            <div v-if="postEditForm.ad_link" style="margin-top: 30px; padding: 20px; background: #f8fafc; border-left: 4px solid #3b82f6; border-radius: 4px;" v-html="parsedAdLink"></div>
          </div>
        </div>
        <div class="modal-footer" style="padding: 16px 24px; border-top: 1px solid #e2e8f0; display: flex; justify-content: flex-end;">
          <button type="button" class="btn-cancel" @click="showPostPreviewModal = false" style="margin: 0;">關閉預覽</button>
        </div>
      </div>
    </div>

    <!-- Hotel CSR Preview Modal Overlay -->
    <div class="modal-overlay" v-if="showHotelPreviewModal" @click.self="showHotelPreviewModal = false">
      <div class="modal-content" style="max-width: 950px; width: 95%; max-height: 85vh; display: flex; flex-direction: column; padding: 0;">
        <div class="modal-header" style="padding: 16px 24px; border-bottom: 1px solid #e2e8f0;">
          <h3 style="margin: 0; font-size: 18px;">旅館前台模擬預覽 (CSR 畫面)</h3>
          <button class="btn-close-modal" @click="showHotelPreviewModal = false">×</button>
        </div>
        <div class="modal-body" style="padding: 30px; overflow-y: auto; background: #f8f9fa; flex: 1;">
          <div style="max-width: 900px; margin: 0 auto; background: white; padding: 30px; border-radius: 12px; box-shadow: 0 5px 20px rgba(0,0,0,0.05);">
            
            <div style="margin-bottom: 25px;">
              <h1 style="font-size: 28px; color: #1e293b; font-weight: 700; margin: 0 0 10px 0;">{{ editForm.name || "未命名旅館" }}</h1>
              <div style="color: #64748b; font-size: 14px;">
                首頁 &gt; {{ editForm.area || "地區" }} &gt; <span style="color: #94a3b8;">{{ editForm.name || "未命名旅館" }}</span>
              </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 30px;">
              <div>
                <div style="background: #f1f5f9; border-radius: 8px; overflow: hidden; height: 300px; display: flex; align-items: center; justify-content: center; position: relative;">
                  <img
                    v-if="editForm.images && editForm.images.length > 0"
                    :src="editForm.images[0]"
                    alt="Hotel Main Image"
                    style="width: 100%; height: 100%; object-fit: cover;"
                  />
                  <div v-else style="color: #94a3b8; font-size: 14px;">暫無圖片</div>
                </div>
                <div style="display: flex; gap: 8px; margin-top: 10px; overflow-x: auto; padding-bottom: 5px;" v-if="editForm.images && editForm.images.length > 1">
                  <div
                    v-for="(img, idx) in editForm.images.slice(1)"
                    :key="idx"
                    style="width: 60px; height: 45px; border-radius: 4px; overflow: hidden; flex-shrink: 0; background: #e2e8f0;"
                  >
                    <img :src="img" style="width: 100%; height: 100%; object-fit: cover;" />
                  </div>
                </div>
              </div>

              <div style="display: flex; flex-direction: column; gap: 20px;">
                <div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 20px; display: flex; flex-direction: column; gap: 12px;">
                  <div v-if="editForm.address" style="display: flex; font-size: 15px;">
                    <span style="color: #64748b; font-weight: 600; width: 60px; flex-shrink: 0;">地址：</span>
                    <span style="color: #334155;">{{ editForm.address }}</span>
                  </div>
                  <div v-if="editForm.phone" style="display: flex; font-size: 15px;">
                    <span style="color: #64748b; font-weight: 600; width: 60px; flex-shrink: 0;">電話：</span>
                    <span style="color: #3b82f6; font-weight: 600;">{{ editForm.phone }}</span>
                  </div>
                  <div v-if="editForm.website" style="display: flex; font-size: 15px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
                    <span style="color: #64748b; font-weight: 600; width: 60px; flex-shrink: 0;">網站：</span>
                    <span style="color: #3b82f6;">{{ editForm.website }}</span>
                  </div>
                  <div v-if="editForm.email" style="display: flex; font-size: 15px;">
                    <span style="color: #64748b; font-weight: 600; width: 60px; flex-shrink: 0;">信箱：</span>
                    <span style="color: #334155;">{{ editForm.email }}</span>
                  </div>
                </div>

                <div style="background: #fdf2f2; border: 1px solid #fde2e2; border-radius: 8px; padding: 20px;">
                  <h3 style="font-size: 16px; color: #9b1c1c; font-weight: 700; margin: 0 0 12px 0; border-bottom: 1px solid #fcd2d2; padding-bottom: 8px;">價格資訊</h3>
                  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
                    <div>
                      <div style="font-size: 13px; color: #7f8c8d;">平日住宿</div>
                      <div style="font-size: 18px; font-weight: 700; color: #e74c3c;" v-if="editForm.pricing.weekday_stay">NT$ {{ editForm.pricing.weekday_stay }}</div>
                      <div style="font-size: 16px; color: #95a5a6; font-weight: 600;" v-else>無資訊</div>
                    </div>
                    <div>
                      <div style="font-size: 13px; color: #7f8c8d;">假日住宿</div>
                      <div style="font-size: 18px; font-weight: 700; color: #e74c3c;" v-if="editForm.pricing.holiday_stay">NT$ {{ editForm.pricing.holiday_stay }}</div>
                      <div style="font-size: 16px; color: #95a5a6; font-weight: 600;" v-else>無資訊</div>
                    </div>
                    <div>
                      <div style="font-size: 13px; color: #7f8c8d;">平日休息</div>
                      <div style="font-size: 16px; font-weight: 700; color: #2c3e50;" v-if="editForm.pricing.weekday_rest">
                        NT$ {{ editForm.pricing.weekday_rest }} <span style="font-size: 12px; font-weight: normal; color: #7f8c8d;">({{ editForm.pricing.weekday_rest_hours }}小時)</span>
                      </div>
                      <div style="font-size: 16px; color: #95a5a6; font-weight: 600;" v-else>無資訊</div>
                    </div>
                    <div>
                      <div style="font-size: 13px; color: #7f8c8d;">假日休息</div>
                      <div style="font-size: 16px; font-weight: 700; color: #2c3e50;" v-if="editForm.pricing.holiday_rest">
                        NT$ {{ editForm.pricing.holiday_rest }} <span style="font-size: 12px; font-weight: normal; color: #7f8c8d;">({{ editForm.pricing.holiday_rest_hours }}小時)</span>
                      </div>
                      <div style="font-size: 16px; color: #95a5a6; font-weight: 600;" v-else>無資訊</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div style="margin-bottom: 30px;">
              <h3 style="font-size: 18px; color: #2c3e50; font-weight: 700; border-bottom: 2px solid #3b82f6; padding-bottom: 8px; margin-bottom: 15px;">入住及退房資訊</h3>
              <p style="font-size: 15px; color: #475569; line-height: 1.6; white-space: pre-line; margin: 0 0 20px 0;" v-if="editForm.stay_info">{{ editForm.stay_info }}</p>
              <p style="font-size: 15px; color: #94a3b8; font-style: italic;" v-else>未設定入住退房資訊</p>

              <h3 style="font-size: 18px; color: #2c3e50; font-weight: 700; border-bottom: 2px solid #3b82f6; padding-bottom: 8px; margin-bottom: 15px;">住宿須知</h3>
              <p style="font-size: 15px; color: #475569; line-height: 1.6; white-space: pre-line; margin: 0;" v-if="editForm.housing_rules">{{ editForm.housing_rules }}</p>
              <p style="font-size: 15px; color: #94a3b8; font-style: italic;" v-else>未設定住宿須知</p>
            </div>

            <div style="background: #f8fafc; padding: 25px; border-radius: 8px; border: 1px dashed #cbd5e1;">
              <h3 style="font-size: 18px; color: #2c3e50; font-weight: 700; margin: 0 0 12px 0;">旅館簡介</h3>
              <p style="font-size: 15px; color: #334155; line-height: 1.7; white-space: pre-line; margin: 0;" v-if="editForm.description">{{ editForm.description }}</p>
              <p style="font-size: 15px; color: #94a3b8; font-style: italic; margin: 0;" v-else>未輸入旅館簡介</p>
            </div>

          </div>
        </div>
        <div class="modal-footer" style="padding: 16px 24px; border-top: 1px solid #e2e8f0; display: flex; justify-content: flex-end;">
          <button type="button" class="btn-cancel" @click="showHotelPreviewModal = false" style="margin: 0;">關閉預覽</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
});
import { ref, onMounted, watch, computed } from "vue";
import { joinURL } from "ufo";
import MarkdownIt from "markdown-it";

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
});

const config = useRuntimeConfig();
const baseURL = config.app.baseURL;
const backendAPI = config.public.backendApiUrl || "http://localhost:8080"; // Base URL of the Go backend

const snackbar = ref({
  show: false,
  message: "",
  type: "success",
});
let snackbarTimeout: any = null;

const showSnackbar = (
  message: string,
  type: "success" | "error" = "success",
) => {
  snackbar.value.message = message;
  snackbar.value.type = type;
  snackbar.value.show = true;
  if (snackbarTimeout) {
    clearTimeout(snackbarTimeout);
  }
  snackbarTimeout = setTimeout(() => {
    snackbar.value.show = false;
  }, 3000);
};

const cities = ref<any[]>([]);
const hotelCategories = ref<any[]>([]);

// Homepage selector states
const homepageSelections = ref<Record<string, string[]>>({
  台北: ["", "", "", "", "", ""],
  新北: ["", "", "", "", "", ""],
  桃園: ["", "", "", "", "", ""],
  台中: ["", "", "", "", "", ""],
  台南: ["", "", "", "", "", ""],
  高雄: ["", "", "", "", "", ""],
});
const allHotelsList = ref<any[]>([]);
const loadingHomepage = ref(false);
const savingHomepage = ref(false);
const activeHomepageCity = ref("台北");
const homepageCities = ["台北", "新北", "桃園", "台中", "台南", "高雄"];
const homepageSearchQuery = ref("");

const selectedCount = computed(() => {
  const city = activeHomepageCity.value;
  return homepageSelections.value[city]
    ? homepageSelections.value[city].filter((id) => id !== "").length
    : 0;
});

const isHotelSelected = (id: string) => {
  const city = activeHomepageCity.value;
  return homepageSelections.value[city]
    ? homepageSelections.value[city].includes(id)
    : false;
};

const getSelectedHotelInSlot = (slotIdx: number) => {
  const city = activeHomepageCity.value;
  const hotelId = homepageSelections.value[city]
    ? homepageSelections.value[city][slotIdx]
    : "";
  if (!hotelId) return null;
  return (
    allHotelsList.value.find((h) => h.id === hotelId) || {
      id: hotelId,
      name: "未知旅館 (已刪除或不存在)",
    }
  );
};

const selectHotelForSlots = (id: string) => {
  const city = activeHomepageCity.value;
  if (!homepageSelections.value[city]) return;
  const current = [...homepageSelections.value[city]];
  const emptyIndex = current.indexOf("");
  if (emptyIndex !== -1) {
    current[emptyIndex] = id;
    homepageSelections.value[city] = current;
  }
};

const removeHotelFromSlot = (idx: number) => {
  const city = activeHomepageCity.value;
  if (!homepageSelections.value[city]) return;
  const current = [...homepageSelections.value[city]];
  current.splice(idx, 1);
  current.push(""); // Pad with empty string to maintain length 6
  homepageSelections.value[city] = current;
};

const moveSlot = (idx: number, direction: number) => {
  const city = activeHomepageCity.value;
  const targetIdx = idx + direction;
  if (targetIdx < 0 || targetIdx > 5 || !homepageSelections.value[city]) return;
  const current = [...homepageSelections.value[city]];
  const temp = current[idx];
  current[idx] = current[targetIdx];
  current[targetIdx] = temp;
  homepageSelections.value[city] = current;
};

const draggedIndex = ref<number | null>(null);

const handleDragStart = (e: DragEvent, idx: number) => {
  draggedIndex.value = idx;
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", idx.toString());
  }
};

const handleDragOver = (e: DragEvent, idx: number) => {
  e.preventDefault();
};

const handleDrop = (e: DragEvent, targetIdx: number) => {
  e.preventDefault();
  const sourceIdx = draggedIndex.value;
  if (sourceIdx === null || sourceIdx === targetIdx) return;

  const city = activeHomepageCity.value;
  if (!homepageSelections.value[city]) return;
  const current = [...homepageSelections.value[city]];

  const temp = current[sourceIdx];
  current[sourceIdx] = current[targetIdx];
  current[targetIdx] = temp;

  homepageSelections.value[city] = current;
};

const handleDragEnd = () => {
  draggedIndex.value = null;
};

const filteredAvailableHotels = computed(() => {
  const city = activeHomepageCity.value;
  const query = homepageSearchQuery.value.trim().toLowerCase();
  let list = allHotelsList.value.filter((h) => h.area === city);
  if (query) {
    list = list.filter(
      (h) =>
        h.id.toLowerCase().includes(query) ||
        h.name.toLowerCase().includes(query),
    );
  }
  return list;
});

// Two-panel Region/City selector state
const regions = ref<any[]>([]);
const selectedCities = ref<string[]>([]);
const showRegionPicker = ref(false);
const activeRegionIdx = ref(0);
const showEditRegionPicker = ref(false);
const activeEditRegionIdx = ref(0);
const isNewHotel = ref(false);

const isDarkMode = ref(false);
const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value;
  if (process.client) {
    localStorage.setItem("cms-theme", isDarkMode.value ? "dark" : "light");
  }
};

const isDeploying = ref(false);
const showDeployConfirmModal = ref(false);
const deployStatus = ref<"idle" | "deploying" | "success" | "error">("idle");
const deployMessage = ref("");

const openDeployModal = () => {
  deployStatus.value = "idle";
  deployMessage.value = "";
  isDeploying.value = false;
  showDeployConfirmModal.value = true;
};

const closeDeployModal = () => {
  if (isDeploying.value) return;
  showDeployConfirmModal.value = false;
};

const executeFrontendDeploy = async () => {
  isDeploying.value = true;
  deployStatus.value = "deploying";
  try {
    const res = await fetch(`${backendAPI}/api/deploy`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token.value}`,
        "Content-Type": "application/json",
      },
    });
    const data = await res.json();
    if (res.ok) {
      deployStatus.value = "success";
      deployMessage.value = data.message || "已成功觸發部署！";
    } else {
      deployStatus.value = "error";
      deployMessage.value = data.error || "未知錯誤";
    }
  } catch (err: any) {
    console.error("Failed to deploy:", err);
    deployStatus.value = "error";
    deployMessage.value = "觸發部署時發生錯誤，請檢查網路連線或伺服器狀態";
  } finally {
    isDeploying.value = false;
  }
};

const showPostPreviewModal = ref(false);
const showHotelPreviewModal = ref(false);

const openPostPreviewModal = () => {
  showPostPreviewModal.value = true;
};

const openHotelPreviewModal = () => {
  showHotelPreviewModal.value = true;
};

const parsedPostContent = computed(() => {
  const content = postEditForm.value.content || "";
  const isHtml = content.includes("<p>") || content.includes("<h3>") || content.includes("<ul") || content.includes("<ol>");
  let rendered = isHtml ? content : md.render(content);
  
  // Replace naked Imgur links converted to anchors by linkify
  rendered = rendered.replace(/<a href="(https:\/\/i\.imgur\.com\/[^"]+)">[^<]+<\/a>/g, (match, url) => {
    return `<img src="${url}" style="width: 100%; max-width: 100%; height: auto; display: block; margin: 20px 0;" />`;
  });
  
  // Replace remaining naked Imgur links not in attributes
  rendered = rendered.replace(/(?<!["'(\/])(https:\/\/i\.imgur\.com\/[a-zA-Z0-9.]+)(?!["'])/g, (match) => {
    return `<img src="${match}" style="width: 100%; max-width: 100%; height: auto; display: block; margin: 20px 0;" />`;
  });
  
  return rendered;
});

const parsedAdLink = computed(() => {
  const ad = postEditForm.value.ad_link || "";
  if (ad.includes("<a") || ad.includes("<p") || ad.includes("<span")) {
    return ad;
  }
  return md.render(ad);
});

// Posts Management States (Admin Only)
const posts = ref<any[]>([]);
const postCurrentPage = ref(1);
const postTotalPages = ref(1);
const postLimit = 10;
const postSearchQuery = ref("");
const loadingPostsList = ref(false);
const selectedPost = ref<any>(null);
const postSaving = ref(false);
const postSuccessMsg = ref("");
const postErrorMsg = ref("");
const newTagText = ref("");
const showPostDeleteConfirmModal = ref(false);
const postToDeleteId = ref<number | null>(null);
const postDeleteConfirmInput = ref("");

interface PostEditForm {
  id?: number;
  title: string;
  tags: string[];
  image: string;
  content: string;
  ad_link: string;
  seo_title: string;
  seo_keywords: string;
  seo_description: string;
}

const postEditForm = ref<PostEditForm>({
  title: "",
  tags: [],
  image: "",
  content: "",
  ad_link: "",
  seo_title: "",
  seo_keywords: "",
  seo_description: "",
});

const addTag = () => {
  const val = newTagText.value.trim();
  if (!val) return;
  if (!postEditForm.value.tags) {
    postEditForm.value.tags = [];
  }
  if (postEditForm.value.tags.length >= 3) {
    alert("分類標籤最多只能設定三個");
    return;
  }
  if (!postEditForm.value.tags.includes(val)) {
    postEditForm.value.tags.push(val);
  }
  newTagText.value = "";
};

const removeTag = (idx: number) => {
  if (postEditForm.value.tags) {
    postEditForm.value.tags.splice(idx, 1);
  }
};

// Auth & Session States
const token = ref("");
const loginUserEmail = ref("");
const userRole = ref("vendor"); // admin or vendor
const activeSection = ref("hotels"); // hotels or users

const loginEmail = ref("");
const loginPassword = ref("");
const loginError = ref("");
const loggingIn = ref(false);

// Sidebar & Hotel List
const hotels = ref<any[]>([]);
const currentPage = ref(1);
const totalPages = ref(1);
const limit = 12;
const searchQuery = ref("");
const selectedArea = ref(""); // kept for backward compat; multi-select uses selectedCities
const loadingList = ref(false);
const quickId = ref("");

// Editing State
const selectedHotel = ref<any>(null);
const saving = ref(false);
const successMsg = ref("");
const errorMsg = ref("");
const hotelEditorTabs = [
  { id: "basic", label: "基本資料" },
  { id: "content", label: "住宿內容" },
  { id: "booking", label: "價格與訂房" },
  { id: "images", label: "圖片管理" },
];
const activeHotelTab = ref("basic");
const setHotelEditorTab = (tabId: string) => {
  activeHotelTab.value = tabId;
};

// Users Management States (Admin Only)
const users = ref<any[]>([]);
const loadingUsers = ref(false);
const showUserModal = ref(false);
const editingUser = ref<any>(null);
const userFormEmail = ref("");
const userFormPassword = ref("");
const userFormRole = ref("vendor");
const userFormError = ref("");
const savingUser = ref(false);
const userSuccessMsg = ref("");
const userErrorMsg = ref("");

interface EditForm {
  id: string;
  name: string;
  area: string;
  address: string;
  phone: string;
  fax: string;
  website: string;
  email: string;
  category: string;
  stay_info: string;
  housing_rules: string;
  pricing: {
    weekday_stay: number;
    holiday_stay: number;
    weekday_rest_hours: number;
    weekday_rest: number;
    holiday_rest_hours: number;
    holiday_rest: number;
  };
  description: string;
  booking_link: string;
  is_disabled: boolean;
  images: string[];
}

const editForm = ref<EditForm>({
  id: "",
  name: "",
  area: "",
  address: "",
  phone: "",
  fax: "",
  website: "",
  email: "",
  category: "",
  stay_info: "",
  housing_rules: "",
  pricing: {
    weekday_stay: 0,
    holiday_stay: 0,
    weekday_rest_hours: 0,
    weekday_rest: 0,
    holiday_rest_hours: 0,
    holiday_rest: 0,
  },
  description: "",
  booking_link: "",
  is_disabled: false,
  images: [],
});

// Fetch city and hotel category options from the unified categories API.
const fetchCategories = async () => {
  try {
    const [citiesRes, hotelCategoriesRes] = await Promise.all([
      fetch(`${backendAPI}/api/categories?type=city`),
      fetch(`${backendAPI}/api/categories?type=hotel_category`),
    ]);

    if (citiesRes.ok) {
      const cityData = await citiesRes.json();
      if (Array.isArray(cityData) && cityData.length > 0) {
        cities.value = cityData.map((category: any) => ({
          id: category.sort_order || category.id,
          name: category.name,
        }));
      }
    }

    if (hotelCategoriesRes.ok) {
      const categoryData = await hotelCategoriesRes.json();
      if (Array.isArray(categoryData)) {
        hotelCategories.value = categoryData;
      }
    }
  } catch (e) {
    console.error("Failed to fetch categories from DB.", e);
  }
};

// Fetch Regions for two-panel selector
const fetchRegions = async () => {
  try {
    const res = await fetch(`${backendAPI}/api/regions`);
    if (res.ok) {
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        regions.value = data;
      }
    }
  } catch (e) {
    console.error("Failed to fetch regions from API.", e);
  }
};

// City filter change handler
const onCityFilterChange = () => {
  selectedArea.value = selectedCities.value.join(",");
  currentPage.value = 1;
  fetchHotels();
};

const clearCityFilter = () => {
  selectedCities.value = [];
  selectedArea.value = "";
  currentPage.value = 1;
  fetchHotels();
};

// Login
const login = async () => {
  loginError.value = "";
  loggingIn.value = true;
  try {
    const res = await fetch(`${backendAPI}/api/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: loginEmail.value,
        password: loginPassword.value,
      }),
    });
    const result = await res.json();
    if (res.ok) {
      token.value = result.token;
      loginUserEmail.value = result.user.email;
      userRole.value = result.user.role || "vendor";
      activeSection.value = "hotels"; // Default to hotels management on login

      localStorage.setItem("admin_token", result.token);
      localStorage.setItem("admin_user", JSON.stringify(result.user));

      fetchHotels();
    } else {
      loginError.value = result.error || "登入失敗";
    }
  } catch (e) {
    loginError.value =
      "登入連線失敗，請確認後端 Go 服務已啟動且資料庫已初始化。";
  } finally {
    loggingIn.value = false;
  }
};

// Logout
const logout = () => {
  token.value = "";
  loginUserEmail.value = "";
  userRole.value = "vendor";
  activeSection.value = "hotels";
  localStorage.removeItem("admin_token");
  localStorage.removeItem("admin_user");
  selectedHotel.value = null;
  hotels.value = [];
  users.value = [];
};

// Fetch Hotels from Go API
const fetchHotels = async () => {
  if (!token.value) return;
  loadingList.value = true;
  try {
    const res = await fetch(
      `${backendAPI}/api/hotels?page=${currentPage.value}&limit=${limit}&area=${encodeURIComponent(selectedArea.value)}&query=${encodeURIComponent(searchQuery.value)}&show_disabled=true`,
    );

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    const result = await res.json();
    if (result.data) {
      hotels.value = result.data;
      totalPages.value = Math.ceil(result.total / limit);
    }
  } catch (e) {
    console.error("Failed to fetch hotels from backend API:", e);
  } finally {
    loadingList.value = false;
  }
};

// Search debounce
let debounceTimer: any = null;
const debounceSearch = () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    currentPage.value = 1;
    fetchHotels();
  }, 300);
};

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    fetchHotels();
  }
};

// Select a hotel to edit
const selectHotel = async (id: string) => {
  if (!id) return;
  isNewHotel.value = false;
  successMsg.value = "";
  errorMsg.value = "";
  activeHotelTab.value = "basic";

  try {
    const res = await fetch(`${backendAPI}/api/hotels/${id}`);

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    if (res.ok) {
      const data = await res.json();
      populateForm(data);
      selectedHotel.value = data;
      return;
    }

    errorMsg.value = `找不到該 ID (${id}) 的旅館資料。`;
  } catch (e) {
    errorMsg.value = "載入旅館資料時發生錯誤：" + e;
  }
};

const populateForm = (data: any) => {
  editForm.value = {
    id: data.id || "",
    name: data.name || "",
    area: data.area || "",
    address: data.address || "",
    phone: data.phone || "",
    fax: data.fax || "",
    website: data.website || "",
    email: data.email || "",
    category: data.category || "",
    stay_info: data.stay_info || "",
    housing_rules: data.housing_rules || "",
    pricing: {
      weekday_stay: Number(data.pricing?.weekday_stay) || 0,
      holiday_stay: Number(data.pricing?.holiday_stay) || 0,
      weekday_rest_hours:
        Number(data.pricing?.weekday_rest_hours ?? data.pricing?.rest_hours) ||
        0,
      weekday_rest: Number(data.pricing?.weekday_rest) || 0,
      holiday_rest_hours:
        Number(data.pricing?.holiday_rest_hours ?? data.pricing?.rest_hours) ||
        0,
      holiday_rest: Number(data.pricing?.holiday_rest) || 0,
    },
    description: data.description || data.stay_info || "",
    booking_link: data.booking_link || "",
    is_disabled: !!data.is_disabled,
    images: data.images || [],
  };
};

// Image URL Management functions
const getImageUrl = (imgName: string) => {
  if (imgName.startsWith("http://") || imgName.startsWith("https://")) {
    return imgName;
  }
  return joinURL(baseURL, `data/images/${imgName}`);
};

const deleteImage = (index: number) => {
  if (editForm.value.images) {
    editForm.value.images.splice(index, 1);
  }
};

const addImageUrl = () => {
  if (!editForm.value.images) {
    editForm.value.images = [];
  }
  editForm.value.images.push("");
};

const draggedImageIndex = ref<number | null>(null);

const handleImageDragStart = (e: DragEvent, idx: number) => {
  draggedImageIndex.value = idx;
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", idx.toString());
  }
};

const handleImageDragOver = (e: DragEvent, idx: number) => {
  e.preventDefault();
};

const handleImageDrop = (e: DragEvent, targetIdx: number) => {
  e.preventDefault();
  const sourceIdx = draggedImageIndex.value;
  if (sourceIdx === null || sourceIdx === targetIdx) return;

  if (!editForm.value.images) return;
  const current = [...editForm.value.images];

  const temp = current[sourceIdx];
  current[sourceIdx] = current[targetIdx];
  current[targetIdx] = temp;

  editForm.value.images = current;
};

const handleImageDragEnd = () => {
  draggedImageIndex.value = null;
};

const getImgurThumbnail = (url: string) => {
  if (!url) return "";
  url = url.trim();
  if (url.startsWith("http://") || url.startsWith("https://")) {
    return url;
  }
  return joinURL(baseURL, `data/images/${url}`);
};

const handleThumbnailError = (e: Event) => {
  const img = e.target as HTMLImageElement;
  img.src = joinURL(baseURL, "data/images/_default.jpg");
};

const isImgurLink = (url: string) => {
  if (!url) return false;
  return url.toLowerCase().includes("imgur.com");
};

const validateAndFormatImgur = (index: number) => {
  let url = editForm.value.images[index].trim();
  if (!url) return;

  if (!url.toLowerCase().includes("imgur.com")) {
    alert("圖片連結僅支援 Imgur 圖床 (imgur.com)！");
    editForm.value.images[index] = "";
    return;
  }

  // Convert imgur page url to direct image link
  if (!url.match(/\.(jpeg|jpg|gif|png|webp)$/i)) {
    const match = url.match(/imgur\.com\/(?:gallery\/|a\/)?([a-zA-Z0-9]+)/);
    if (match && match[1]) {
      url = `https://i.imgur.com/${match[1]}.jpg`;
    } else {
      alert("請輸入有效的 Imgur 圖片網址！");
      editForm.value.images[index] = "";
      return;
    }
  }

  editForm.value.images[index] = url;
};

const selectEditArea = (city: string) => {
  editForm.value.area = city;
  showEditRegionPicker.value = false;
};

const previewImage = (url: string) => {
  if (url && url.trim()) {
    window.open(url.trim(), "_blank");
  }
};

const initNewHotel = () => {
  isNewHotel.value = true;
  selectedHotel.value = { id: "" };
  editForm.value = {
    id: "",
    name: "未命名新旅館",
    area: "",
    address: "",
    phone: "",
    fax: "",
    website: "",
    email: "",
    category: "",
    stay_info: "",
    housing_rules: "",
    pricing: {
      weekday_stay: 0,
      holiday_stay: 0,
      weekday_rest_hours: 0,
      weekday_rest: 0,
      holiday_rest_hours: 0,
      holiday_rest: 0,
    },
    description: "",
    booking_link: "",
    is_disabled: false,
    images: [],
  };
  activeHotelTab.value = "basic";
  successMsg.value = "";
  errorMsg.value = "";
};

const handleImgPreviewError = (e: Event) => {
  const img = e.target as HTMLImageElement;
  img.src = joinURL(baseURL, "data/images/_default.jpg");
};

// Save Hotel Details
const saveHotel = async () => {
  successMsg.value = "";
  errorMsg.value = "";

  // Validate images before save
  if (editForm.value.images && editForm.value.images.length > 0) {
    for (let i = 0; i < editForm.value.images.length; i++) {
      const url = editForm.value.images[i].trim();
      if (!url) continue;
      if (!url.toLowerCase().includes("imgur.com")) {
        errorMsg.value = `第 ${i + 1} 張圖片連結非 Imgur 網址，請確認！`;
        return;
      }
    }
    // Clean up empty images
    editForm.value.images = editForm.value.images.filter(
      (img) => img.trim() !== "",
    );
  }

  saving.value = true;

  try {
    const res = await fetch(`${backendAPI}/api/hotels/${editForm.value.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify(editForm.value),
    });

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    const result = await res.json();
    if (res.ok) {
      successMsg.value = "旅館資料成功儲存到資料庫！";
      selectedHotel.value = result;
      isNewHotel.value = false;
      fetchHotels();
    } else {
      errorMsg.value = result.error || "儲存資料失敗。";
    }
  } catch (e) {
    errorMsg.value = "儲存連線失敗：" + e;
  } finally {
    saving.value = false;
  }
};

// ==================== Users Account Management CRUD (Admin Only) ====================

const fetchUsers = async () => {
  if (userRole.value !== "admin") return;
  loadingUsers.value = true;
  userSuccessMsg.value = "";
  userErrorMsg.value = "";

  try {
    const res = await fetch(`${backendAPI}/api/users`, {
      headers: {
        Authorization: `Bearer ${token.value}`,
      },
    });

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    const result = await res.json();
    if (res.ok) {
      users.value = result;
    } else {
      userErrorMsg.value = result.error || "無法取得帳號清單";
    }
  } catch (e) {
    userErrorMsg.value = "讀取帳號連線失敗：" + e;
  } finally {
    loadingUsers.value = false;
  }
};

const openCreateUserModal = () => {
  editingUser.value = null;
  userFormEmail.value = "";
  userFormPassword.value = "";
  userFormRole.value = "vendor";
  userFormError.value = "";
  showUserModal.value = true;
};

const openEditUserModal = (user: any) => {
  editingUser.value = user;
  userFormEmail.value = user.email;
  userFormPassword.value = ""; // Clear password field, only change if user enters new password
  userFormRole.value = user.role;
  userFormError.value = "";
  showUserModal.value = true;
};

const saveUser = async () => {
  userFormError.value = "";
  savingUser.value = true;

  const url = editingUser.value
    ? `${backendAPI}/api/users/${editingUser.value.id}`
    : `${backendAPI}/api/users`;

  const method = editingUser.value ? "PUT" : "POST";

  try {
    const res = await fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({
        email: userFormEmail.value,
        password: userFormPassword.value,
        role: userFormRole.value,
      }),
    });

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    const result = await res.json();
    if (res.ok) {
      userSuccessMsg.value = editingUser.value
        ? "帳號設定修改成功！"
        : "帳號新增成功！";
      showUserModal.value = false;
      fetchUsers();
    } else {
      userFormError.value = result.error || "儲存帳號失敗";
    }
  } catch (e) {
    userFormError.value = "儲存帳號連線失敗：" + e;
  } finally {
    savingUser.value = false;
  }
};

const deleteUser = async (id: number) => {
  if (!confirm("您確定要刪除這個管理帳號嗎？刪除後將無法恢復。")) return;

  userSuccessMsg.value = "";
  userErrorMsg.value = "";

  try {
    const res = await fetch(`${backendAPI}/api/users/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token.value}`,
      },
    });

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    const result = await res.json();
    if (res.ok) {
      userSuccessMsg.value = "帳號已成功刪除！";
      fetchUsers();
    } else {
      userErrorMsg.value = result.error || "刪除帳號失敗";
    }
  } catch (e) {
    userErrorMsg.value = "刪除帳號連線失敗：" + e;
  }
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return "-";
  try {
    const d = new Date(dateStr);
    return d.toLocaleString("zh-TW", { hour12: false });
  } catch (e) {
    return dateStr;
  }
};

// ==================== Posts CRUD Management (Admin Only) ====================

const fetchPosts = async () => {
  if (!token.value) return;
  loadingPostsList.value = true;
  try {
    const res = await fetch(
      `${backendAPI}/api/posts?page=${postCurrentPage.value}&limit=${postLimit}&search=${encodeURIComponent(postSearchQuery.value)}`,
    );
    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }
    if (res.ok) {
      const result = await res.json();
      posts.value = result.data || [];
      postTotalPages.value = Math.ceil(result.total / postLimit);
    }
  } catch (e) {
    console.error("Failed to fetch posts:", e);
  } finally {
    loadingPostsList.value = false;
  }
};

const selectPost = async (id: number) => {
  postSuccessMsg.value = "";
  postErrorMsg.value = "";
  try {
    const res = await fetch(`${backendAPI}/api/posts/${id}`);
    if (res.ok) {
      const data = await res.json();
      selectedPost.value = data;
      postEditForm.value = {
        id: data.id,
        title: data.title || "",
        tags: data.tags || [],
        image: data.image || "",
        content: data.content || "",
        ad_link: data.ad_link || "",
        seo_title: data.seo_title || "",
        seo_keywords: data.seo_keywords || "",
        seo_description: data.seo_description || "",
      };
    }
  } catch (e) {
    postErrorMsg.value = "讀取文章失敗: " + e;
  }
};

const openCreatePost = async () => {
  postSuccessMsg.value = "";
  postErrorMsg.value = "";
  selectedPost.value = { id: 0 };
  postEditForm.value = {
    title: "",
    tags: [],
    image: "",
    content: "",
    ad_link: "",
    seo_title: "",
    seo_keywords: "",
    seo_description: "",
  };
};

const savePost = async () => {
  postSaving.value = true;
  postSuccessMsg.value = "";
  postErrorMsg.value = "";

  const isNew = !postEditForm.value.id;
  const url = isNew
    ? `${backendAPI}/api/posts`
    : `${backendAPI}/api/posts/${postEditForm.value.id}`;
  const method = isNew ? "POST" : "PUT";

  try {
    const res = await fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify(postEditForm.value),
    });

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    const result = await res.json();
    if (res.ok) {
      postSuccessMsg.value = "文章儲存成功！";
      selectedPost.value = result;
      postEditForm.value.id = result.id;
      fetchPosts();
    } else {
      postErrorMsg.value = result.error || "儲存文章失敗";
    }
  } catch (e) {
    postErrorMsg.value = "儲存文章連線失敗: " + e;
  } finally {
    postSaving.value = false;
  }
};

const deletePost = (id: number) => {
  postToDeleteId.value = id;
  postDeleteConfirmInput.value = "";
  showPostDeleteConfirmModal.value = true;
};

const executeDeletePost = async () => {
  if (postDeleteConfirmInput.value !== "確認刪除") return;
  const id = postToDeleteId.value;
  if (id === null) return;

  showPostDeleteConfirmModal.value = false;
  postSuccessMsg.value = "";
  postErrorMsg.value = "";
  try {
    const res = await fetch(`${backendAPI}/api/posts/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token.value}`,
      },
    });
    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }
    if (res.ok) {
      postSuccessMsg.value = "文章已成功刪除";
      selectedPost.value = null;
      fetchPosts();
    } else {
      const result = await res.json();
      postErrorMsg.value = result.error || "刪除文章失敗";
    }
  } catch (e) {
    postErrorMsg.value = "刪除文章連線失敗: " + e;
  }
};

const fetchHomepageHotels = async () => {
  loadingHomepage.value = true;
  try {
    const hotelsRes = await fetch(
      `${backendAPI}/api/hotels?limit=10000&show_disabled=true`,
    );
    if (hotelsRes.ok) {
      const data = await hotelsRes.json();
      allHotelsList.value = data.data || [];
    }

    const selectionsRes = await fetch(`${backendAPI}/api/homepage-hotels`);
    if (selectionsRes.ok) {
      const selectedList = await selectionsRes.json();
      homepageSelections.value = {
        台北: ["", "", "", "", "", ""],
        新北: ["", "", "", "", "", ""],
        桃園: ["", "", "", "", "", ""],
        台中: ["", "", "", "", "", ""],
        台南: ["", "", "", "", "", ""],
        高雄: ["", "", "", "", "", ""],
      };

      for (const item of selectedList) {
        if (
          homepageSelections.value[item.city] &&
          item.sort_order >= 1 &&
          item.sort_order <= 6
        ) {
          homepageSelections.value[item.city][item.sort_order - 1] =
            item.hotel_id;
        }
      }
    }
  } catch (e) {
    errorMsg.value = "載入首頁選定資料失敗：" + e;
  } finally {
    loadingHomepage.value = false;
  }
};

const saveHomepageHotels = async () => {
  savingHomepage.value = true;
  const selectionsList: any[] = [];
  for (const city of homepageCities) {
    const slots = homepageSelections.value[city];
    for (let i = 0; i < slots.length; i++) {
      const hotelId = slots[i];
      if (hotelId) {
        selectionsList.push({
          city: city,
          sort_order: i + 1,
          hotel_id: hotelId,
        });
      }
    }
  }

  try {
    const res = await fetch(`${backendAPI}/api/homepage-hotels`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({ selections: selectionsList }),
    });

    if (res.status === 401) {
      logout();
      loginError.value = "登入逾期，請重新登入。";
      return;
    }

    if (res.ok) {
      successMsg.value = "首頁精選旅館設定儲存成功！";
    } else {
      const data = await res.json();
      errorMsg.value = data.error || "儲存首頁精選旅館失敗";
    }
  } catch (e) {
    errorMsg.value = "連線失敗：" + e;
  } finally {
    savingHomepage.value = false;
  }
};

const getHotelsForCity = (city: string) => {
  return allHotelsList.value.filter((h) => h.area === city);
};

watch(successMsg, (val) => {
  if (val) {
    showSnackbar(val, "success");
    setTimeout(() => {
      successMsg.value = "";
    }, 3000);
  }
});
watch(errorMsg, (val) => {
  if (val) {
    showSnackbar(val, "error");
    setTimeout(() => {
      errorMsg.value = "";
    }, 3000);
  }
});
watch(postSuccessMsg, (val) => {
  if (val) {
    showSnackbar(val, "success");
    setTimeout(() => {
      postSuccessMsg.value = "";
    }, 3000);
  }
});
watch(postErrorMsg, (val) => {
  if (val) {
    showSnackbar(val, "error");
    setTimeout(() => {
      postErrorMsg.value = "";
    }, 3000);
  }
});
watch(userSuccessMsg, (val) => {
  if (val) {
    showSnackbar(val, "success");
    setTimeout(() => {
      userSuccessMsg.value = "";
    }, 3000);
  }
});
watch(userErrorMsg, (val) => {
  if (val) {
    showSnackbar(val, "error");
    setTimeout(() => {
      userErrorMsg.value = "";
    }, 3000);
  }
});
watch(loginError, (val) => {
  if (val) {
    showSnackbar(val, "error");
    setTimeout(() => {
      loginError.value = "";
    }, 3000);
  }
});

onMounted(() => {
  if (process.client) {
    const savedTheme = localStorage.getItem("cms-theme");
    if (savedTheme === "dark") {
      isDarkMode.value = true;
    }

    // Intercept fetch for token refreshing and authorization error handling (kick out)
    const originalFetch = window.fetch;
    window.fetch = async (input, init) => {
      const res = await originalFetch(input, init);
      
      // Auto-extract refresh token from X-Refresh-Token header if available
      if (res.headers && typeof res.headers.get === "function") {
        const refreshToken = res.headers.get("x-refresh-token") || res.headers.get("X-Refresh-Token");
        if (refreshToken) {
          token.value = refreshToken;
          localStorage.setItem("admin_token", refreshToken);
        }
      }
      
      // If 401 Unauthorized or 403 Forbidden is returned, kick the user out to login screen
      if (res.status === 401 || res.status === 403) {
        // Skip for the login page itself to avoid alerting on incorrect credentials
        const isLoginUrl = typeof input === "string" && input.includes("/api/auth/login");
        if (!isLoginUrl) {
          token.value = "";
          userRole.value = "vendor";
          localStorage.removeItem("admin_token");
          localStorage.removeItem("admin_user");
          alert("您的登入狀態已過期或權限不足，系統已自動將您登出。");
        }
      }
      return res;
    };
  }
  token.value = localStorage.getItem("admin_token") || "";

  const userJson = localStorage.getItem("admin_user");
  if (userJson) {
    const user = JSON.parse(userJson);
    loginUserEmail.value = user.email || "";
    userRole.value = user.role || "vendor";
  }

  fetchCategories();
  fetchRegions();
  if (token.value) {
    fetchHotels();
  }
});
</script>

<style scoped>
/* Reset & Base Variables */
.admin-page {
  height: 100vh;
  max-height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background-color: #f8fafc;
  color: #1e293b;
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
    Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  transition:
    background-color 0.3s,
    color 0.3s;
}

/* Custom Webkit Scrollbars for independent scroll zones */
.admin-sidebar::-webkit-scrollbar,
.admin-main::-webkit-scrollbar,
.users-workspace::-webkit-scrollbar,
.rich-editor::-webkit-scrollbar,
textarea::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.admin-sidebar::-webkit-scrollbar-track,
.admin-main::-webkit-scrollbar-track,
.users-workspace::-webkit-scrollbar-track {
  background: transparent;
}

.admin-sidebar::-webkit-scrollbar-thumb,
.admin-main::-webkit-scrollbar-thumb,
.users-workspace::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.admin-sidebar::-webkit-scrollbar-thumb:hover,
.admin-main::-webkit-scrollbar-thumb:hover,
.users-workspace::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* ------------------- LIGHT MODE (DEFAULT) ------------------- */
.admin-header {
  height: 56px;
  min-height: 56px;
  background-color: #ffffff;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  padding: 0 24px;
  color: #1e293b;
  box-sizing: border-box;
  transition: all 0.3s;
}

.header-container {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.admin-title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.025em;
  color: #0f172a;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* Theme Toggle Button */
.btn-theme-toggle {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #475569;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn-theme-toggle:hover {
  background-color: #f1f5f9;
  color: #0f172a;
}

.theme-icon {
  width: 16px;
  height: 16px;
}

.user-badge {
  font-size: 13px;
  color: #475569;
  background-color: #f1f5f9;
  padding: 6px 12px;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.user-icon {
  width: 14px;
  height: 14px;
  color: #3b82f6;
}

.btn-logout {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #475569;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-logout:hover {
  background-color: #fee2e2;
  border-color: #fca5a5;
  color: #ef4444;
}

.btn-home {
  background-color: #3b82f6;
  color: white;
  text-decoration: none;
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  transition: background-color 0.2s;
}

.btn-home:hover {
  background-color: #2563eb;
}

.btn-deploy {
  background-color: #3b82f6; /* Blue background */
  color: white;
  border: none;
  padding: 10px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: background-color 0.2s;
  width: 100%;
  margin-bottom: 10px;
}

.btn-deploy:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn-deploy:disabled {
  background-color: #94a3b8;
  cursor: not-allowed;
  opacity: 0.6;
}

/* Spinner and Animation for Deploying state */
.deploy-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin-deploy 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin-deploy {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* ------------------- MAIN WRAPPER ------------------- */
.admin-main-wrapper {
  flex: 1;
  display: flex;
  overflow: hidden;
  height: calc(100vh - 56px);
}

/* ------------------- LEFT NAVIGATION (LIGHT MENU) ------------------- */
.admin-nav {
  width: 220px;
  flex-shrink: 0;
  background-color: #f8fafc;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  height: 100%;
  transition: all 0.3s;
}

.nav-brand {
  padding: 24px;
  font-size: 12px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #64748b;
}

.nav-links {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 12px;
  gap: 6px;
}

.nav-link {
  background: transparent;
  border: none;
  color: #475569;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  text-align: left;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.2s;
}

.nav-link:hover {
  background-color: #e2e8f0;
  color: #0f172a;
}

.nav-link.active {
  background-color: #eff6ff;
  color: #2563eb;
}

.nav-icon {
  width: 18px;
  height: 18px;
}

.nav-footer {
  padding: 20px 16px;
  border-top: 1px solid #e2e8f0;
}

.nav-footer .btn-theme-toggle {
  width: 100%;
  justify-content: center;
  margin-bottom: 10px;
  padding: 10px;
}

.hotel-editor-tabs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  padding: 14px 24px;
  background: #ffffff;
}

.hotel-editor-tabs button {
  padding: 11px 14px;
  border: 0;
  border-radius: 8px;
  background: #e2e8f0;
  color: #475569;
  cursor: pointer;
  text-align: center;
  transition: 0.2s ease;
}

.hotel-editor-tabs button:hover {
  background: #cbd5e1;
  color: #1e293b;
}

.hotel-editor-tabs button.active {
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.22);
}

.hotel-editor-tabs span {
  font-size: 13px;
  font-weight: 700;
}

.pricing-settings {
  display: grid;
  gap: 18px;
}

.pricing-section {
  padding: 18px;
  border-radius: 12px;
  background: #f8fafc;
}

.pricing-section-heading {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.pricing-section-heading h4,
.pricing-section-heading p {
  margin: 0;
}

.pricing-section-heading h4 {
  color: #1e293b;
  font-size: 15px;
}

.pricing-section-heading p {
  margin-top: 3px;
  color: #94a3b8;
  font-size: 12px;
}

.pricing-section-icon {
  display: grid;
  width: 36px;
  height: 36px;
  place-items: center;
  border-radius: 10px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 13px;
  font-weight: 800;
}

.pricing-section-icon.rest {
  background: #dcfce7;
  color: #15803d;
}

.pricing-option-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.pricing-option {
  position: relative;
  padding: 18px 14px 2px;
  border-radius: 10px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.06);
}

.pricing-option-label {
  position: absolute;
  top: 0;
  left: 14px;
  padding: 3px 9px;
  border-radius: 0 0 7px 7px;
  background: #2563eb;
  color: #ffffff;
  font-size: 11px;
  font-weight: 700;
}

.pricing-option-label.holiday {
  background: #e11d48;
}

.pricing-field-pair {
  display: grid;
  grid-template-columns: minmax(0, 0.8fr) minmax(0, 1.2fr);
  gap: 10px;
  padding-top: 10px;
}

@media (max-width: 760px) {
  .pricing-option-grid,
  .pricing-field-pair {
    grid-template-columns: 1fr;
  }
}

.btn-home-nav {
  display: block;
  text-align: center;
  background-color: #e2e8f0;
  color: #475569;
  text-decoration: none;
  padding: 10px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-home-nav:hover {
  background-color: #cbd5e1;
  color: #0f172a;
}

/* ------------------- WORKSPACE CONTAINER ------------------- */

.hotels-workspace,
.posts-workspace {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
}
.admin-workspace {
  flex: 1;
  height: 100%;
  overflow: hidden;
  display: flex;
}

/* ------------------- SIDEBAR (LIST) ------------------- */
.admin-sidebar {
  width: 300px;
  flex-shrink: 0;
  height: 100%;
  border-right: 1px solid #e2e8f0;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  transition: all 0.3s;
}

.search-box {
  padding: 20px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-search {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
  background-color: white;
  color: #0f172a;
  box-sizing: border-box;
  transition: all 0.2s;
}

.input-search:focus {
  outline: none;
  border-color: #3b82f6;
}

.btn-create-hotel {
  background-color: #10b981;
  color: white;
  border: none;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.2s;
  width: 100%;
}

.btn-create-hotel:hover {
  background-color: #059669;
}

/* ---- Two-Panel Region/City Selector ---- */
.region-selector-wrapper {
  position: relative;
  width: 100%;
}

.btn-region-trigger {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
  background-color: white;
  color: #0f172a;
  box-sizing: border-box;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  text-align: left;
  transition: all 0.2s;
}

.btn-region-trigger:hover {
  border-color: #3b82f6;
}

.filter-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: #64748b;
}

.btn-region-trigger span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chevron-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: #94a3b8;
  transition: transform 0.2s;
}

.chevron-icon.open {
  transform: rotate(180deg);
}

.region-picker-popover {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  width: 360px;
  max-width: calc(100vw - 40px);
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow:
    0 10px 25px -5px rgba(0, 0, 0, 0.1),
    0 8px 10px -6px rgba(0, 0, 0, 0.05);
  z-index: 200;
  overflow: hidden;
}

.region-picker-body {
  display: flex;
  min-height: 240px;
  max-height: 320px;
}

.region-list {
  width: 110px;
  flex-shrink: 0;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  background-color: #f8fafc;
}

.region-item {
  background: transparent;
  border: none;
  padding: 14px 16px;
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  text-align: left;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition: all 0.15s;
}

.region-item:hover {
  background-color: #eff6ff;
  color: #2563eb;
}

.region-item.active {
  background-color: white;
  color: #2563eb;
  border-left-color: #2563eb;
}

.city-list {
  flex: 1;
  padding: 12px 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.city-checkbox {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 6px;
  font-size: 14px;
  color: #334155;
  cursor: pointer;
  transition: background-color 0.15s;
}

.city-checkbox:hover {
  background-color: #f1f5f9;
}

.city-picker-item {
  width: 100%;
  text-align: left;
  background: transparent;
  border: none;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 14px;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s;
}

.city-picker-item:hover {
  background-color: #f1f5f9;
  color: #2563eb;
}

.city-picker-item.active {
  background-color: #eff6ff;
  color: #2563eb;
  font-weight: 600;
}

.city-checkbox input[type="checkbox"],
.city-checkbox input[type="radio"] {
  width: 16px;
  height: 16px;
  accent-color: #3b82f6;
  cursor: pointer;
}

.region-picker-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-top: 1px solid #e2e8f0;
  background-color: #fafbfc;
}

.btn-clear-filter {
  background: transparent;
  border: none;
  color: #64748b;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: color 0.2s;
}

.btn-clear-filter:hover {
  color: #2563eb;
}

.btn-close-picker {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 8px 24px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-close-picker:hover {
  background-color: #2563eb;
}

/* Loading & Empty state */
.loading-state,
.empty-state {
  padding: 40px;
  text-align: center;
  color: #64748b;
  font-size: 14px;
}

/* Hotel / Post lists */
.hotel-list,
.post-list {
  list-style: none;
  padding: 0;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.hotel-item,
.post-item {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
  cursor: pointer;
  transition: background-color 0.2s;
}

.hotel-item:hover,
.post-item:hover {
  background-color: #f8fafc;
}

.hotel-item.active,
.post-item.active {
  background-color: #eff6ff;
  border-left: 4px solid #3b82f6;
}

.hotel-item-name {
  font-size: 14px;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 6px;
  line-height: 1.4;
}

.hotel-item-info {
  display: flex;
  gap: 8px;
}

.badge {
  background-color: #f1f5f9;
  color: #475569;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.category {
  background-color: #dbeafe;
  color: #1e40af;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

/* Post list details */
.post-item-title {
  font-size: 14px;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 6px;
  line-height: 1.4;
}

.post-item-info {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag-badge {
  background-color: #f1f5f9;
  color: #475569;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
}

/* Pagination */
.pagination {
  padding: 12px 20px;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination button {
  background: white;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  background-color: #f1f5f9;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination span {
  font-size: 12px;
  color: #64748b;
  font-weight: 600;
}

/* ------------------- EDITOR AREA ------------------- */
.admin-main {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  padding: 32px;
  background-color: #f8fafc;
  box-sizing: border-box;
  transition: all 0.3s;
}

.no-selection {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #94a3b8;
  text-align: center;
}

.no-selection .icon {
  margin-bottom: 16px;
}

.no-selection h2 {
  font-size: 20px;
  font-weight: 700;
  color: #475569;
  margin: 0 0 8px 0;
}

.no-selection p {
  margin: 0 0 20px 0;
  font-size: 14px;
}

.quick-load {
  display: flex;
  gap: 10px;
  max-width: 320px;
  width: 100%;
}

.quick-load input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 13px;
}

.quick-load button {
  padding: 8px 16px;
  background-color: #475569;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
}

.quick-load button:hover {
  background-color: #334155;
}

/* Editor Wrapper */
.editor-wrapper {
  max-width: 800px;
  margin: 0 auto;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 16px;
}

.editor-header h2 {
  font-size: 22px;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.id-tag {
  font-size: 13px;
  font-weight: 600;
  color: #94a3b8;
  background: #f1f5f9;
  padding: 2px 8px;
  border-radius: 4px;
}

.btn-save {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-left: 12px;
}

.btn-save:hover {
  background-color: #2563eb;
}

.btn-preview {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background-color: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #334155;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
  margin-left: 12px;
}

.btn-preview:hover {
  background-color: #e2e8f0;
  color: #0f172a;
}

/* Vertical Form Cards */
.editor-vertical-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

.card h3 {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  border-left: 4px solid #3b82f6;
  padding-left: 10px;
}

/* Fields Styling */
.form-group {
  margin-bottom: 18px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  margin-bottom: 6px;
}

.form-group input[type="text"],
.form-group input[type="email"],
.form-group input[type="password"],
.form-group textarea,
.select-field {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 14px;
  box-sizing: border-box;
  color: #0f172a;
  background-color: white;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group textarea:focus,
.select-field:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-group textarea {
  resize: vertical;
}

/* Grids */
.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 18px;
}

.form-grid-3 {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
  margin-bottom: 18px;
}

@media (max-width: 600px) {
  .hotel-editor-tabs {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    padding: 10px 12px;
  }

  .form-grid-2,
  .form-grid-3 {
    grid-template-columns: 1fr;
  }
}

/* Rich Text Toolbars */
.editor-toolbar {
  background-color: #f8fafc;
  border: 1px solid #cbd5e1;
  border-bottom: none;
  padding: 8px;
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.editor-toolbar button {
  background: white;
  border: 1px solid #e2e8f0;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.editor-toolbar button:hover {
  background-color: #f1f5f9;
  border-color: #cbd5e1;
}

.btn-clear-format {
  margin-left: auto;
  color: #ef4444 !important;
}

.rich-editor {
  border: 1px solid #cbd5e1;
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
  min-height: 180px;
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
  background-color: white;
  font-size: 14px;
  line-height: 1.6;
}

.rich-editor:focus {
  outline: none;
  border-color: #3b82f6;
}

/* Images UI */
.image-editor-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.no-images {
  text-align: center;
  color: #64748b;
  padding: 20px;
  font-size: 13px;
  background-color: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
}

.image-urls-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.image-url-row {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 12px;
  padding: 16px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background-color: #f8fafc;
  transition: all 0.2s ease;
}

.image-url-row.dragging {
  opacity: 0.5;
  border-color: #3b82f6;
  background-color: #eff6ff;
  border-style: solid;
  transform: scale(0.98);
}

.image-url-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.img-order-badge {
  background-color: #e2e8f0;
  color: #475569;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.input-image-url {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 14px;
  background-color: white;
  color: #1e293b;
  transition: all 0.2s;
}

.input-image-url:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.image-url-row .row-actions {
  display: flex;
  gap: 6px;
}

.btn-row-action {
  background: white;
  border: 1px solid #cbd5e1;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
}

.btn-row-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-row-delete {
  background-color: #ef4444;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-row-delete:hover {
  background-color: #dc2626;
}

.btn-row-preview {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-row-preview:hover {
  background-color: #2563eb;
}

.btn-row-preview:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.required {
  color: #ef4444;
  margin-left: 2px;
}

.url-validation-warning {
  width: 100%;
  color: #ef4444;
  font-size: 12px;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.url-preview {
  width: 100%;
  margin-top: 10px;
}

.url-preview img {
  max-width: 150px;
  max-height: 100px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid #cbd5e1;
}

.btn-add-img {
  width: 100%;
  padding: 10px;
  background-color: #f1f5f9;
  color: #475569;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-add-img:hover {
  background-color: #e2e8f0;
  color: #0f172a;
}

/* ------------------- USERS MANAGEMENT ------------------- */
.users-workspace {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  padding: 32px;
  background-color: #f8fafc;
  box-sizing: border-box;
  transition: all 0.3s;
}

.workspace-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 16px;
}

.workspace-header h2 {
  font-size: 22px;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.btn-create-user {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-create-user:hover {
  background-color: #2563eb;
}

.users-table-wrapper {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 14px;
}

.users-table th,
.users-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.users-table th {
  background-color: #f8fafc;
  font-weight: 700;
  color: #475569;
}

.role-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
}

.role-badge.admin {
  background-color: #fee2e2;
  color: #991b1b;
}

.role-badge.vendor {
  background-color: #e0f2fe;
  color: #075985;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.btn-table-edit {
  background: white;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-table-edit:hover {
  background-color: #f1f5f9;
  border-color: #cbd5e1;
}

.btn-table-delete {
  background-color: #fee2e2;
  color: #ef4444;
  border: 1px solid #fca5a5;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-table-delete:hover {
  background-color: #fca5a5;
  color: #b91c1c;
}

/* ------------------- MODAL OVERLAY ------------------- */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 460px;
  max-width: 90%;
  padding: 32px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  box-sizing: border-box;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  color: #0f172a;
}

.btn-close-modal {
  background: transparent;
  border: none;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
}

.btn-close-modal:hover {
  color: #475569;
}

.select-role {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 14px;
  background-color: white;
  box-sizing: border-box;
  cursor: pointer;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
}

.btn-cancel {
  background-color: white;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
}

.btn-cancel:hover {
  background-color: #f8fafc;
}

.btn-submit {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
}

.btn-submit:hover {
  background-color: #2563eb;
}

/* ------------------- POSTS WORKSPACE ------------------- */
.btn-create-post {
  width: 100%;
  padding: 10px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-create-post:hover {
  background-color: #2563eb;
}

.btn-delete-post {
  background-color: #ef4444;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-delete-post:hover {
  background-color: #dc2626;
}

.tags-input-container {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 16px;
  background-color: #f8fafc;
}

.tags-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.tag-chip {
  background-color: #3b82f6;
  color: white;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-remove-tag {
  background: transparent;
  border: none;
  color: white;
  cursor: pointer;
  font-size: 14px;
}

.tag-input-row {
  display: flex;
  gap: 10px;
}

.tag-input-row input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 14px;
  background-color: white;
  color: #1e293b;
  transition: all 0.2s;
  box-sizing: border-box;
}

.tag-input-row input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.btn-add-tag {
  padding: 10px 18px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

/* Warnings and Alerts */
.alert-success {
  background-color: #d1fae5;
  color: #065f46;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 20px;
  border: 1px solid #a7f3d0;
}

.alert-danger {
  background-color: #fee2e2;
  color: #991b1b;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 20px;
  border: 1px solid #fca5a5;
}

.warning-icon {
  width: 14px;
  height: 14px;
  color: #ef4444;
}

/* ------------------- LOGIN CARD ------------------- */
.login-wrapper {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f1f5f9;
}

.login-card {
  width: 400px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.05);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-icon svg {
  width: 40px;
  height: 40px;
  color: #2563eb;
}

.login-header h2 {
  margin: 12px 0 6px 0;
  color: #0f172a;
  font-size: 24px;
  font-weight: 700;
}

.login-header p {
  margin: 0;
  color: #64748b;
  font-size: 14px;
}

.login-card label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  margin-bottom: 8px;
}

.login-card input {
  width: 100%;
  padding: 10px 14px;
  border-radius: 8px;
  background-color: #ffffff;
  border: 1px solid #cbd5e1;
  color: #0f172a;
  font-size: 14px;
  box-sizing: border-box;
  transition: border-color 0.2s;
}

.login-card input:focus {
  outline: none;
  border-color: #3b82f6;
}

.btn-login {
  width: 100%;
  padding: 12px;
  background-color: #3b82f6;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-weight: 700;
  font-size: 15px;
  cursor: pointer;
  margin-top: 24px;
  transition: background-color 0.2s;
}

.btn-login:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  color: #64748b;
}

/* ------------------- DARK MODE CLASS OVERRIDES ------------------- */
.dark-mode {
  background-color: #0f172a;
  color: #cbd5e1;
}

.dark-mode .admin-header {
  background-color: #0f172a;
  border-bottom: 1px solid #1e293b;
  color: #f8fafc;
}

.dark-mode .admin-title {
  color: #f8fafc;
}

.dark-mode .btn-theme-toggle {
  border-color: #334155;
  color: #94a3b8;
}

.dark-mode .btn-theme-toggle:hover {
  background-color: #1e293b;
  color: #f8fafc;
}

.dark-mode .user-badge {
  color: #94a3b8;
  background-color: #1e293b;
}

.dark-mode .user-icon {
  color: #38bdf8;
}

.dark-mode .btn-logout {
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .btn-logout:hover {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
  color: #f87171;
}

.dark-mode .btn-preview {
  background-color: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .btn-preview:hover {
  background-color: #334155;
  color: #ffffff;
}

.dark-mode .admin-nav {
  background-color: #0f172a;
  border-right: 1px solid #1e293b;
}

.dark-mode .nav-link {
  color: #94a3b8;
}

.dark-mode .nav-link:hover {
  background-color: #1e293b;
  color: #ffffff;
}

.dark-mode .nav-link.active {
  background-color: #1e293b;
  color: #38bdf8;
}

.dark-mode .btn-home-nav {
  background-color: #1e293b;
  color: #cbd5e1;
}

.dark-mode .btn-home-nav:hover {
  background-color: #334155;
  color: white;
}

.dark-mode .hotel-editor-tabs {
  background: #0f172a;
}

.dark-mode .hotel-editor-tabs button {
  background: #1e293b;
  color: #94a3b8;
}

.dark-mode .hotel-editor-tabs button:hover {
  background: #334155;
  color: #e2e8f0;
}

.dark-mode .hotel-editor-tabs button.active {
  background: #0284c7;
  color: #ffffff;
  box-shadow: 0 4px 10px rgba(2, 132, 199, 0.25);
}

.dark-mode .pricing-section {
  background: #0f172a;
}

.dark-mode .pricing-section-heading h4 {
  color: #cbd5e1;
}

.dark-mode .pricing-option {
  background: #1e293b;
  box-shadow: none;
}

.dark-mode :deep(.integer-input-unit) {
  background: #334155;
  color: #cbd5e1;
}

.dark-mode :deep(.wang-editor-shell) {
  border-color: #334155;
}

.dark-mode :deep(.wang-editor-toolbar),
.dark-mode :deep(.wang-editor-content) {
  background: #0f172a;
  color: #cbd5e1;
}

.dark-mode .admin-sidebar {
  background-color: #0f172a;
  border-right: 1px solid #1e293b;
}

.dark-mode .search-box {
  border-bottom: 1px solid #1e293b;
}

.dark-mode .input-search {
  background-color: #0b0f19;
  border-color: #334155;
  color: #f8fafc;
}

.dark-mode .input-search:focus {
  border-color: #38bdf8;
}

/* Dark Mode: Region Picker */
.dark-mode .btn-region-trigger {
  background-color: #0b0f19;
  border-color: #334155;
  color: #f8fafc;
}

.dark-mode .btn-region-trigger:hover {
  border-color: #38bdf8;
}

.dark-mode .filter-icon {
  color: #94a3b8;
}

.dark-mode .chevron-icon {
  color: #64748b;
}

.dark-mode .region-picker-popover {
  background-color: #1e293b;
  border-color: #334155;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
}

.dark-mode .region-list {
  background-color: #0f172a;
  border-right-color: #334155;
}

.dark-mode .region-item {
  color: #94a3b8;
}

.dark-mode .region-item:hover {
  background-color: #1e293b;
  color: #38bdf8;
}

.dark-mode .region-item.active {
  background-color: #1e293b;
  color: #38bdf8;
  border-left-color: #38bdf8;
}

.dark-mode .city-checkbox {
  color: #cbd5e1;
}

.dark-mode .city-checkbox:hover {
  background-color: #334155;
}

.dark-mode .region-picker-footer {
  border-top-color: #334155;
  background-color: #0f172a;
}

.dark-mode .btn-clear-filter {
  color: #94a3b8;
}

.dark-mode .btn-clear-filter:hover {
  color: #38bdf8;
}

.dark-mode .btn-close-picker {
  background-color: #38bdf8;
  color: #0f172a;
}

.dark-mode .btn-close-picker:hover {
  background-color: #0ea5e9;
}

.dark-mode .hotel-item,
.dark-mode .post-item {
  border-bottom: 1px solid #1e293b;
}

.dark-mode .hotel-item:hover,
.dark-mode .post-item:hover {
  background-color: #1e293b;
}

.dark-mode .hotel-item.active,
.dark-mode .post-item.active {
  background-color: #1e293b;
  border-left-color: #38bdf8;
}

.dark-mode .hotel-item-name,
.dark-mode .post-item-title {
  color: #f8fafc;
}

.dark-mode .badge {
  background-color: #1e293b;
  color: #94a3b8;
}

.dark-mode .category {
  background-color: rgba(56, 189, 248, 0.1);
  color: #38bdf8;
}

.dark-mode .tag-badge {
  background-color: #1e293b;
  color: #94a3b8;
}

.dark-mode .pagination {
  border-top: 1px solid #1e293b;
}

.dark-mode .pagination button {
  background-color: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .pagination button:hover:not(:disabled) {
  background-color: #334155;
}

.dark-mode .admin-main {
  background-color: #0b0f19;
}

.dark-mode .no-selection h2 {
  color: #cbd5e1;
}

.dark-mode .quick-load input {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .editor-header {
  border-bottom: 1px solid #1e293b;
}

.dark-mode .editor-header h2 {
  color: #ffffff;
}

.dark-mode .id-tag {
  background: #1e293b;
  color: #cbd5e1;
}

.dark-mode .card {
  background-color: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2);
}

.dark-mode .card h3 {
  color: #ffffff;
  border-left-color: #38bdf8;
}

.dark-mode .form-group label {
  color: #94a3b8;
}

.dark-mode .form-group input[type="text"],
.dark-mode .form-group input[type="email"],
.dark-mode .form-group input[type="password"],
.dark-mode .form-group textarea,
.dark-mode .select-field,
.dark-mode .rich-editor {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .form-group input:focus,
.dark-mode .form-group textarea:focus,
.dark-mode .select-field:focus,
.dark-mode .rich-editor:focus {
  border-color: #38bdf8;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.15);
}

.dark-mode .editor-toolbar {
  background-color: #0f172a;
  border-color: #334155;
}

.dark-mode .editor-toolbar button {
  background-color: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .editor-toolbar button:hover {
  background-color: #334155;
  border-color: #475569;
}

.dark-mode .image-url-row {
  border-color: #334155;
  background-color: #0f172a;
}

.dark-mode .image-url-row.dragging {
  border-color: #38bdf8;
  background-color: rgba(56, 189, 248, 0.1);
}

.dark-mode .img-order-badge {
  background-color: #334155;
  color: #cbd5e1;
}

.dark-mode .btn-row-action {
  background-color: #1e293b;
  border-color: #334155;
  color: white;
}

.dark-mode .url-preview img {
  border-color: #334155;
}

.dark-mode .btn-add-img {
  background-color: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .btn-add-img:hover {
  background-color: #334155;
  color: white;
}

.dark-mode .input-image-url {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .input-image-url:focus {
  border-color: #38bdf8;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.15);
}

.dark-mode .city-checkbox input[type="radio"] {
  accent-color: #38bdf8;
}

.dark-mode .users-workspace {
  background-color: #0b0f19;
}

.dark-mode .workspace-header {
  border-bottom: 1px solid #1e293b;
}

.dark-mode .workspace-header h2 {
  color: #ffffff;
}

.dark-mode .users-table-wrapper {
  border-color: #334155;
}

.dark-mode .users-table th {
  background-color: #1e293b;
  color: #cbd5e1;
}

.dark-mode .users-table td {
  border-bottom-color: #1e293b;
  color: #cbd5e1;
}

.dark-mode .btn-table-edit {
  background-color: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .btn-table-edit:hover {
  background-color: #334155;
}

.dark-mode .modal-content {
  background-color: #1e293b;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .modal-header h3 {
  color: #ffffff;
}

.dark-mode .btn-cancel {
  background-color: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}

.dark-mode .btn-cancel:hover {
  background-color: #334155;
}

.dark-mode .tags-input-container {
  border-color: #334155;
  background-color: #0f172a;
}

.dark-mode .tag-input-row input {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .tag-input-row input:focus {
  border-color: #38bdf8;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.15);
}

.dark-mode .btn-add-tag {
  background-color: #38bdf8;
  color: #0f172a;
}

.dark-mode .btn-add-tag:hover {
  background-color: #0ea5e9;
}

.dark-mode .login-wrapper {
  background-color: #0f172a;
}

.dark-mode .login-card {
  background-color: #1e293b;
  border-color: #334155;
}

.dark-mode .login-header h2 {
  color: #ffffff;
}

.dark-mode .login-card label {
  color: #cbd5e1;
}

.dark-mode .login-card input {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .btn-login {
  background-color: #38bdf8;
  color: #0f172a;
}

.dark-mode .btn-login:hover:not(:disabled) {
  background-color: #0ea5e9;
}

.dark-mode .btn-row-preview {
  background-color: #2563eb;
}

.dark-mode .btn-row-preview:hover:not(:disabled) {
  background-color: #1d4ed8;
}

.dark-mode .btn-create-hotel {
  background-color: #059669;
}

.dark-mode .btn-create-hotel:hover {
  background-color: #047857;
}

.dark-mode .city-picker-item {
  color: #cbd5e1;
}

.dark-mode .city-picker-item:hover {
  background-color: #334155;
  color: #38bdf8;
}

.dark-mode .city-picker-item.active {
  background-color: #1e293b;
  color: #38bdf8;
}

/* Checkbox Label */
.checkbox-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #334155;
  font-weight: 500;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #ef4444;
  cursor: pointer;
}

.dark-mode .checkbox-label {
  color: #cbd5e1;
}

/* Snackbar styling */
.snackbar {
  position: fixed;
  top: 24px;
  right: 24px;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 20px;
  border-radius: 10px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  transform: translateY(-20px);
  opacity: 0;
  pointer-events: none;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  min-width: 280px;
  max-width: 420px;
}

.snackbar.show {
  transform: translateY(0);
  opacity: 1;
  pointer-events: auto;
}

.snackbar.success {
  background-color: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #065f46;
}

.snackbar.error {
  background-color: #fef2f2;
  border: 1px solid #fca5a5;
  color: #991b1b;
}

.snackbar-content {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  font-weight: 500;
}

.snackbar-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.snackbar-close {
  background: transparent;
  border: none;
  font-size: 18px;
  font-weight: 600;
  color: currentColor;
  cursor: pointer;
  opacity: 0.7;
  padding: 0;
  line-height: 1;
}

.snackbar-close:hover {
  opacity: 1;
}

/* Dark Mode Snackbar overrides */
.dark-mode .snackbar.success {
  background-color: #064e3b;
  border-color: #047857;
  color: #ecfdf5;
}

.dark-mode .snackbar.error {
  background-color: #7f1d1d;
  border-color: #b91c1c;
  color: #fef2f2;
}

/* Homepage Editor Styles */
.homepage-workspace {
  flex: 1;
  padding: 30px;
  background-color: #f8fafc;
  overflow-y: auto;
}

.homepage-editor-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.city-tabs {
  display: flex;
  gap: 10px;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
}

.city-tab-btn {
  padding: 10px 20px;
  background: transparent;
  border: none;
  font-size: 16px;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  border-bottom: 3px solid transparent;
  margin-bottom: -11px;
  transition: all 0.2s;
}

.city-tab-btn:hover {
  color: #3b82f6;
}

.city-tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

/* Transfer Container Layout */
.transfer-container {
  display: flex;
  gap: 24px;
  align-items: stretch;
  margin-top: 15px;
}

.transfer-panel {
  flex: 1;
  background-color: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  height: 600px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.panel-header-custom {
  padding: 16px 20px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background-color: #f8fafc;
  border-radius: 12px 12px 0 0;
}

.panel-header-custom h3 {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.search-wrapper-custom {
  width: 100%;
}

.input-search-inline {
  width: 100%;
  padding: 8px 12px;
  font-size: 14px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  outline: none;
  transition: all 0.2s;
  box-sizing: border-box;
}

.input-search-inline:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.panel-body {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.scrollable-y {
  overflow-y: auto;
}

/* Items in Left list */
.transfer-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #f1f5f9;
  transition: background-color 0.2s;
}

.transfer-item:hover {
  background-color: #f8fafc;
}

.item-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-id {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  background-color: #f1f5f9;
  padding: 2px 6px;
  border-radius: 4px;
  align-self: flex-start;
}

.item-name {
  font-size: 14px;
  color: #1e293b;
  font-weight: 600;
}

.item-addr {
  font-size: 12px;
  color: #64748b;
}

.btn-select-item {
  padding: 6px 12px;
  font-size: 13px;
  font-weight: 600;
  color: white;
  background-color: #3b82f6;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-select-item:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn-select-item:disabled {
  background-color: #e2e8f0;
  color: #94a3b8;
  cursor: not-allowed;
}

.empty-state-text {
  text-align: center;
  color: #94a3b8;
  padding: 40px 0;
  font-size: 14px;
}

/* Status Badge */
.status-badge {
  font-size: 12px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 4px;
  background-color: #fee2e2;
  color: #ef4444;
  align-self: flex-start;
}

.status-badge.valid {
  background-color: #d1fae5;
  color: #065f46;
}

/* Right Panel Slot Row */
.slot-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  transition: all 0.2s ease;
}

.slot-row.dragging {
  opacity: 0.5;
  transform: scale(0.98);
}

.slot-row.dragging .slot-content {
  border-color: #3b82f6;
  background-color: #eff6ff;
  border-style: solid;
}

.slot-number-badge {
  font-size: 12px;
  font-weight: 700;
  color: #475569;
  background-color: #e2e8f0;
  width: 50px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  flex-shrink: 0;
}

.slot-content {
  flex: 1;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
  min-height: 48px;
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background-color: #f8fafc;
}

.selected-hotel-display {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.display-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.display-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.slot-actions {
  display: flex;
  gap: 6px;
}

/* Drag handle style */
.drag-handle {
  transition: color 0.2s;
}
.drag-handle:hover {
  color: #3b82f6 !important;
}

.btn-slot-remove {
  background-color: #fee2e2;
  border: 1px solid #fca5a5;
  color: #ef4444;
  font-weight: 600;
  border-radius: 4px;
  padding: 2px 8px;
  font-size: 12px;
  cursor: pointer;
}

.btn-slot-remove:hover {
  background-color: #fca5a5;
  color: white;
}

.empty-slot-placeholder {
  color: #94a3b8;
  font-size: 13px;
  font-style: italic;
}

/* Dark Mode Shuttle Box Styles */
.dark-mode .transfer-panel {
  background-color: #1e293b;
  border-color: #334155;
}

.dark-mode .panel-header-custom {
  background-color: #0f172a;
  border-bottom-color: #334155;
}

.dark-mode .panel-header-custom h3 {
  color: #ffffff;
}

.dark-mode .input-search-inline {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

.dark-mode .transfer-item {
  border-bottom-color: #334155;
}

.dark-mode .transfer-item:hover {
  background-color: #0f172a;
}

.dark-mode .item-id {
  background-color: #0f172a;
  color: #94a3b8;
}

.dark-mode .item-name {
  color: #ffffff;
}

.dark-mode .item-addr {
  color: #94a3b8;
}

.dark-mode .btn-select-item:disabled {
  background-color: #334155;
  color: #64748b;
}

.dark-mode .slot-number-badge {
  color: #cbd5e1;
  background-color: #334155;
}

.dark-mode .slot-content {
  border-color: #475569;
  background-color: #0f172a;
}

.dark-mode .display-name {
  color: #ffffff;
}

.dark-mode .slot-row.dragging .slot-content {
  border-color: #38bdf8;
  background-color: rgba(56, 189, 248, 0.1);
}

.dark-mode .drag-handle:hover {
  color: #38bdf8 !important;
}

.dark-mode .btn-slot-remove {
  background-color: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.3);
  color: #fca5a5;
}

.dark-mode .btn-slot-remove:hover {
  background-color: #ef4444;
  color: white;
}
</style>
